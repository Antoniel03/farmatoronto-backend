-- falta empleados.cedula y empleados.fecha_ingreso
SELECT empleados.id,empleados.nombre,empleados.apellido,empleados.cedula,
      empleados.direccion,empleados.telefono, usuarios.correo, empleados.fecha_nacimiento,
      empleados.cargo

FROM empleados
JOIN usuarios on usuarios.codempleado = empleados.id


