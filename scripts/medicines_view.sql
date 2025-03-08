SELECT 
    a.id AS medicamento_id,
    a.nombre AS medicamento_nombre,
    a.componenteprincipal,
    a.accion,
    a.presentacion,
    a.precio,
    s.cantidad AS stock_cantidad,
    f.nombre AS farmacia_nombre
FROM 
    medicamentos a
INNER JOIN 
    Medic_labs b ON a.id = b.codmedicamento
INNER JOIN 
    stock s ON a.id = s.medicamento_id    
    INNER JOIN 
    farmacia_sucursal f ON s.farmacia_id = f.id 
WHERE 
    a.nombre = "?";
