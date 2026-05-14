import { Warehouse } from '../warehouses/warehouse.entity';
import { Product } from '../products/product.entity';
export declare class Inventory {
    id: number;
    warehouseId: number;
    productId: number;
    quantity: number;
    lastUpdated: Date;
    warehouse: Warehouse;
    product: Product;
}
