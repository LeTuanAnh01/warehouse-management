import { Repository } from 'typeorm';
import { Product } from './product.entity';
export declare class ProductsService {
    private productsRepository;
    constructor(productsRepository: Repository<Product>);
    findAll(): Promise<Product[]>;
    findOne(id: number): Promise<Product>;
    create(dto: Partial<Product>): Promise<Product>;
    update(id: number, dto: Partial<Product>): Promise<Product>;
    remove(id: number): Promise<{
        message: string;
    }>;
}
