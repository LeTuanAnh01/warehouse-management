import { Repository } from 'typeorm';
import { Inventory } from './inventory.entity';
export declare class InventoryService {
    private inventoryRepository;
    constructor(inventoryRepository: Repository<Inventory>);
    findAll(): Promise<Inventory[]>;
    findByWarehouse(warehouseId: number): Promise<Inventory[]>;
    findOne(id: number): Promise<Inventory>;
    create(dto: Partial<Inventory>): Promise<Inventory>;
    update(id: number, dto: Partial<Inventory>): Promise<Inventory>;
    remove(id: number): Promise<{
        message: string;
    }>;
}
