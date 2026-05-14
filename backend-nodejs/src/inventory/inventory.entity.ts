import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  UpdateDateColumn,
  ManyToOne,
  JoinColumn,
} from 'typeorm';
import { Warehouse } from '../warehouses/warehouse.entity';
import { Product } from '../products/product.entity';

@Entity('inventory')
export class Inventory {
  @PrimaryGeneratedColumn()
  id: number;

  @Column({ name: 'warehouse_id' })
  warehouseId: number;

  @Column({ name: 'product_id' })
  productId: number;

  @Column({ default: 0 })
  quantity: number;

  @UpdateDateColumn({ name: 'last_updated' })
  lastUpdated: Date;

  @ManyToOne(() => Warehouse, { eager: true })
  @JoinColumn({ name: 'warehouse_id' })
  warehouse: Warehouse;

  @ManyToOne(() => Product, { eager: true })
  @JoinColumn({ name: 'product_id' })
  product: Product;
}
