-- falta empleados.cedula y empleados.fecha_ingreso
SELECT empleados.id,empleados.nombre,empleados.apellido, empleados.fecha_nacimiento,
      empleados.direccion, empleados.telefono, usuarios.correo
FROM empleados
JOIN usuarios on usuarios.codempleado = empleados.id


