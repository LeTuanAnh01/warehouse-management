import { ProductsService } from './products.service';
export declare class ProductsController {
    private productsService;
    constructor(productsService: ProductsService);
    findAll(): Promise<import("./product.entity").Product[]>;
    findOne(id: string): Promise<import("./product.entity").Product>;
    create(dto: any): Promise<import("./product.entity").Product>;
    update(id: string, dto: any): Promise<import("./product.entity").Product>;
    remove(id: string): Promise<{
        message: string;
    }>;
}
