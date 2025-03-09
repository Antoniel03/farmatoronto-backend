SELECT COUNT(*) FROM medicamentos JOIN medic_labs ON medic_labs.codmedicamento = medicamentos.id
    JOIN laboratorio ON laboratorio.id = medic_labs.codlaboratorio
    JOIN stock ON stock.medicamento_id = medicamentos.id
    JOIN farmacia_sucursal ON farmacia_sucursal.id = stock.farmacia_id
    JOIN accion_terapeutica ON accion_terapeutica.id = medicamentos.accion_id
    JOIN medic_monodrogas ON medic_monodrogas.codmedicamento = medicamentos.id
    JOIN monodrogas ON monodrogas.id = medic_monodrogas.codmonodroga
    JOIN ciudad ON ciudad.id = farmacia_sucursal.ciudad_id

