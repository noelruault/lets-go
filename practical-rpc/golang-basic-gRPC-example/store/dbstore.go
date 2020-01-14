package store

/*
	CREATE TABLE films (
	id             serial PRIMARY KEY,
	title          varchar(255) NOT NULL,
	director       varchar(255) NOT NULL,
	producer       varchar(255) NOT NULL,
	release_date   date NOT NULL DEFAULT NOW()::date
	);

	INSERT INTO films (id, title, director, producer, release_date) VALUES
	('4', 'A New Hope', 'George Lucas', 'Gary Kurtz, Rick McCallum', '1977-5-25' ),
	('5', 'The Empire Strikes Back', 'Irvin Kershner', 'Gary Kurtz, Rick McCallum', '1980-5-17' ),
	('6', 'Return of the Jedi', 'Richard Marquand', 'Howard G. Kazanjian, George Lucas, Rick McCallum', '1983-5-25' );
	```

	-- ALTER TABLE films ADD PRIMARY KEY (id);
*/
