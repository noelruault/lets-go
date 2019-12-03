# Practical gRPC

## Chapter 9. Load balancing

### Client load balancing

Client load balancing adds more complexity to the architecture.
The client is much more sophisticated, since it needs to apply strategies to equilibrate the traffic.

- Thick client: The client will need to keep track on the health of each backend, in order to redirect traffic elsewhere if a backend malfunction is detected.
- Lookaside Load Balancer: Communicates the client that is the best backend server to communicate with.

### Proxy load balancing

The main advantage of proxy load balancing is the simplicity of the client. It will only need a single endpoint to create a connection with. All of the workload problems, the security issues and the awareness of the health of every backend server, will be completely transparent to the client.

- Transport Level: The proxy just checks that the socket is open, in order to know if the backend is up and running. There is no payload treatment. Client data is just copied to the backend connection.
- Application Level: On the proxy, the HTTP/2 protocol is parsed in order to inspect each request and make decisions on the fly.

## Chapter 10. Service evolution with gRPC

Details required to maintain a functioning service over time while making changes.
Related: [byte-level details about protobuf encoding for specific types](https://developers.google.com/protocol-buffers/docs/encoding)

### Binary and source compatibility

When using the codegen tools for gRPC (protoc, protoc-gen-grpc, and related plugins) there are two interfaces you need to be concerned with: the ABI and the API.

- The application binary interface (ABI): Refers to the binary interface of the tool that’s being used to generate the protobuf code as well as any code that the generated code will rely on.
- The application programming interface (API): Refers to the signatures of the methods we, as the consumer, can call the objects that are available for us, etc.

To avoid the need for breaking changes in the future, there are a number of things you can do right at the beginning when you define the service:

- Versioning a Service.
- Define custom request and response objects.

```go
syntax = "proto3";

package practical_grpc.v1;

service MyService {
  rpc MyMethod(MyMethodRequest) returns (MyMethodResponse);
}

message MyMethodRequest {
}
```

### Maintaining wire compatibility

In gRPC, calls are only distinguishable by their service and method names. For example, given the following protobuf, the target for a call would be practical_grpc.v1.SomeService/SomeMethod.

```go
Protocol Buffer

syntax = "proto3";

package practical_grpc.v1;

service SomeService {
  rpc SomeMethod(SomeMethodRequest) returns (stream SomeMethodResponse);
}
// ...
```

Notice how the method’s request and response types are not included here. These, along with the cardinality of the method signature (unary vs streaming) are implicit in the call.

This means, at least in theory, that you could change a method from unary to streaming (by prepending a message with stream), or change the name of a message type, assuming it had the same structure. However, as discussed in the previous section on ABI/API compatibility, these changes would likely introduce breaking changes for your existing clients, and should be avoided.

#### Implicit versus explicit values on the wire

Message names in gRPC calls are implicit, Field names are also implicit, except when they’re referenced directly by the JSON name or a FieldMask.

Field tags on the other hand are explicit.

Tags cannot be changed without breaking existing clients, and their reuse should be avoided.

Some types are equivalent on the wire. For example, int32 and int64 are the same on the wire. This doesn’t mean they’re treated equally by the client libraries in all languages.

#### Behavior compatibility

##### Default values for new fields

Default values for new fields: For a full list of the default values for each field type, see the [protobuf docs](https://developers.google.com/protocol-buffers/docs/proto3#default)

There are times, however, when you might need to check for presence and not just the default value. If it’s necessary, you can box primitive types so that you can check for presence. This is accomplished by using wrappers.proto from the Well-Known Types. There’s a wrapper for each of the defined primitive values.

```go
syntax = "proto3";
import "google/protobuf/wrapper.proto";

package practical_grpc.v1;

message SomeMessage {
  // allows us to check for `nil` instead of just the default value
  google.protobuf.StringValue value = 1;
}
```

This relies on the fact that the default value for messages is nil. This is equivalent to boxing in other programming languages you may have seen, where primitive types are wrapped by an object so you can do nil checks.

##### Field masks

**proto**

```go
syntax = "proto3";
import "google/protobuf/field_mask.proto";
// ...
message UpdateThingRequest {
  Thing thing = 1;

  // allow callers to specify which fields should be updated
  google.protobuf.FieldMask mask = 2;
}
```

**client**

```go
// ...
 request = Request.new(
    thing: thing,
    mask: Google::Protobuf::FieldMask.new(paths: %w[thing.id thing.name])
  )
// ...
```

**server**

```go
func (s *MyService) UpdateThing(ctx context.Context, r *pb.UpdateThingRequest) (*pb.UpdateThingResponse, error) {
    if r.Mask != nil {
        // Update only fields in r.Mask.Paths
        log.Print(r.Mask.Paths)
    }

    return &pb.UpdateThingResponse{Thing: r.Thing}, nil
}
```

While this approach does solve the problem with field clearing, it requires you to use field masks everywhere upfront. Once you implement an RPC method, you would need to include a field mask for each request object to protect against future changes.

##### Boxing

Another option is to box the new value. This allows you to do a presence check and only update the value when it’s supplied.

```go
syntax = "proto3";
import "google/protobuf/wrappers.proto";
// ...
message Thing {
  Int64 id    = 1;
  String name = 2;

  // wrap for presence checking
  google.protobuf.BoolValue cool = 3;
}
```

The fields themselves also need wrapping/unwrapping on the client side, which could prove to be unnecessarily verbose.

#### Dealing with errors
