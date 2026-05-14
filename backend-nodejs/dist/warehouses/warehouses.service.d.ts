import { Repository } from 'typeorm';
import { Warehouse } from './warehouse.entity';
export declare class WarehousesService {
    private warehousesRepository;
    constructor(warehousesRepository: Repository<Warehouse>);
    findAll(): Promise<Warehouse[]>;
    findOne(id: number): Promise<Warehouse>;
    create(dto: Partial<Warehouse>): Promise<Warehouse>;
    update(id: number, dto: Partial<Warehouse>): Promise<Warehouse>;
    remove(id: number): Promise<{
        message: string;
    }>;
}
