SELECT medicamentos.id, medicamentos.nombre, medicamentos.componenteprincipal,
       medicamentos.presentacion, laboratorio.nombre, medicamentos.precio, stock.cantidad,
       medicamentos.accion
FROM medicamentos
JOIN Medic_labs ON Medic_labs.codmedicamento = medicamentos.id
JOIN laboratorio ON Medic_labs.codlaboratorio = laboratorio.id
JOIN  stock ON stock.medicamento_id= medicamentos.id
JOIN farmacia_sucursal ON farmacia_sucursal.id=stock.farmacia_id


