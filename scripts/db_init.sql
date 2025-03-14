CREATE TABLE IF NOT EXISTS usuarios(
            id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
            correo VARCHAR(50) UNIQUE,
            contrasena VARCHAR(50),
            codempleado INTEGER UNIQUE,
            FOREIGN KEY(codempleado) REFERENCES empleados(id)
      );

CREATE TABLE IF NOT EXISTS empleados(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,-- id puede ser cedula o pasaporte
    nombre VARCHAR(50),
    apellido VARCHAR(50),
    cargo VARCHAR(50),
    fecha_nacimiento VARCHAR(50),
    direccion VARCHAR(100),
    telefono VARCHAR(12)UNIQUE,
    cedula VARCHAR(12)UNIQUE
);

CREATE TABLE IF NOT EXISTS administrador (
    id INTEGER NOT NULL PRIMARY KEY,
    FOREIGN KEY (id) REFERENCES empleados(id)
);

CREATE TABLE IF NOT EXISTS auxiliar (
    id INTEGER NOT NULL PRIMARY KEY,
    FOREIGN KEY (id) REFERENCES empleados(id)
);

CREATE TABLE IF NOT EXISTS analista (
    id INTEGER NOT NULL PRIMARY KEY,
    FOREIGN KEY (id) REFERENCES empleados(id)
);

CREATE TABLE IF NOT EXISTS pasantes (
    id INTEGER NOT NULL PRIMARY KEY,
    FOREIGN KEY (id) REFERENCES empleados(id)
);


CREATE TABLE IF NOT EXISTS farmaceutico (
    id INTEGER NOT NULL PRIMARY KEY,
    FOREIGN KEY (id) REFERENCES empleados(id)
);


CREATE TABLE IF NOT EXISTS medicamentos(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    nombre VARCHAR(50),
    componenteprincipal VARCHAR(50),
    accion_id INTEGER,
    presentacion VARCHAR(50),
    precio REAL,
    FOREIGN KEY(accion_id) REFERENCES accion_terapeutica(id)
);

CREATE TABLE IF NOT EXISTS accion_terapeutica(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    descripcion VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS monodrogas(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    nombre VARCHAR(50) UNIQUE,
    descripcion VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS medic_monodrogas(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    codmedicamento INTEGER,
    codmonodroga INTEGER,
    FOREIGN KEY(codmedicamento) REFERENCES medicamentos(id),
    FOREIGN KEY(codmonodroga) REFERENCES monodrogas(id)
);

CREATE TABLE IF NOT EXISTS stock(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
    farmacia_id INTEGER,
    medicamento_id INTEGER,
    cantidad INTEGER,
    FOREIGN KEY(farmacia_id) REFERENCES farmacia_sucursal(id),
    FOREIGN KEY(medicamento_id) REFERENCES medicamentos(id)
);

CREATE TABLE IF NOT EXISTS rotacion (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    empleado_id INTEGER,
    sucursal_id INTEGER,
    fecha_inicio DATE,
    fecha_final DATE,
    observaciones VARCHAR(100),
    FOREIGN KEY(empleado_id) REFERENCES empleados(id),
    FOREIGN KEY(sucursal_id) REFERENCES farmacia_sucursal(id)
);

CREATE TABLE IF NOT EXISTS farmacia_sucursal(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    ciudad_id INTEGER UNIQUE,
    direccion VARCHAR(100)UNIQUE,
    telefono VARCHAR(12)UNIQUE,
    FOREIGN KEY(ciudad_id) REFERENCES ciudad(id)
);

CREATE TABLE IF NOT EXISTS ciudad(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    nombre VARCHAR(50)UNIQUE,
    estado VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS compra(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    ordencompra_id INTEGER,
    fechaRecepcion DATE,
    total REAL,
    FOREIGN KEY(ordencompra_id) REFERENCES ordencompra(id)
);

CREATE TABLE IF NOT EXISTS compra_medicamentos(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    compra_id INTEGER,
    medicamento_id INTEGER,
    cantidad INTEGER,
    FOREIGN KEY(compra_id) REFERENCES compra(id),
    FOREIGN KEY(medicamento_id) REFERENCES medicamentos(id)
);

CREATE TABLE IF NOT EXISTS ordencompra(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    codfarmacia INTEGER, 
    codanalista INTEGER,
    codlaboratorio INTEGER,
    fechaEmision DATE,
    formaPago VARCHAR(50),
    FOREIGN KEY(codfarmacia) REFERENCES farmacia_sucursal(id),
    FOREIGN KEY(codanalista) REFERENCES analista(id),
    FOREIGN KEY(codlaboratorio) REFERENCES laboratorio(id)
);

CREATE TABLE IF NOT EXISTS laboratorio(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    nombre VARCHAR(50),
    direccion VARCHAR(100)UNIQUE,
    telefono VARCHAR(12) UNIQUE
);

CREATE TABLE IF NOT EXISTS medic_labs(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    codmedicamento INTEGER,
    codlaboratorio INTEGER,
    FOREIGN KEY(codmedicamento) REFERENCES medicamentos(id),
    FOREIGN KEY(codlaboratorio) REFERENCES laboratorio(id)
);
