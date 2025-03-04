CREATE TABLE IF NOT EXISTS empleados(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			nombre VARCHAR(50),
			apellido VARCHAR(50),
			fecha_nacimiento VARCHAR(50),
			direccion VARCHAR(100),
			telefono VARCHAR(12),
			email VARCHAR(50)
		  );


CREATE TABLE IF NOT EXISTS medicamentos(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			nombre VARCHAR(50),
			componenteprincipal VARCHAR(50),
			precio REAL
		  );


CREATE TABLE IF NOT EXISTS usuarios(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			correo VARCHAR(50) UNIQUE,
			contrasena VARCHAR(50),
      tipo_usuario VARCHAR(15)
      );




