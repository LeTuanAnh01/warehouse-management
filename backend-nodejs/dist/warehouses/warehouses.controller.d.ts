import { WarehousesService } from './warehouses.service';
export declare class WarehousesController {
    private warehousesService;
    constructor(warehousesService: WarehousesService);
    findAll(): Promise<import("./warehouse.entity").Warehouse[]>;
    findOne(id: string): Promise<import("./warehouse.entity").Warehouse>;
    create(dto: any): Promise<import("./warehouse.entity").Warehouse>;
    update(id: string, dto: any): Promise<import("./warehouse.entity").Warehouse>;
    remove(id: string): Promise<{
        message: string;
    }>;
}
