import { InventoryService } from './inventory.service';
export declare class InventoryController {
    private inventoryService;
    constructor(inventoryService: InventoryService);
    findAll(): Promise<import("./inventory.entity").Inventory[]>;
    findByWarehouse(warehouseId: string): Promise<import("./inventory.entity").Inventory[]>;
    findOne(id: string): Promise<import("./inventory.entity").Inventory>;
    create(dto: any): Promise<import("./inventory.entity").Inventory>;
    update(id: string, dto: any): Promise<import("./inventory.entity").Inventory>;
    remove(id: string): Promise<{
        message: string;
    }>;
}
