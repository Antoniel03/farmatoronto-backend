SELECT COUNT(*) FROM medicamentos JOIN medic_labs ON medic_labs.medicamento_id = medicamentos.id
    JOIN laboratorio ON laboratorio.id = medic_labs.laboratorio_id
    JOIN stock ON stock.medicamento_id = medicamentos.id
    JOIN sucursal ON sucursal.id = stock.farmacia_id
    JOIN accion_terapeutica ON accion_terapeutica.id = medicamentos.accion_id
    JOIN ciudad ON ciudad.id = sucursal.ciudad_id
   JOIN medic_monodrogas ON medic_monodrogas.medicamento_id = medicamentos.id
   JOIN monodrogas ON monodrogas.id = medic_monodrogas.monodroga_id


