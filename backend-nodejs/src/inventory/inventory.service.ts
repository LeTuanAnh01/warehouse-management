import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Inventory } from './inventory.entity';

@Injectable()
export class InventoryService {
  constructor(
    @InjectRepository(Inventory)
    private inventoryRepository: Repository<Inventory>,
  ) {}

  findAll() {
    return this.inventoryRepository.find();
  }

  findByWarehouse(warehouseId: number) {
    return this.inventoryRepository.find({ where: { warehouseId } });
  }

  async findOne(id: number) {
    const item = await this.inventoryRepository.findOne({ where: { id } });
    if (!item) throw new NotFoundException('Inventory item not found');
    return item;
  }

  create(dto: Partial<Inventory>) {
    const item = this.inventoryRepository.create(dto);
    return this.inventoryRepository.save(item);
  }

  async update(id: number, dto: Partial<Inventory>) {
    await this.inventoryRepository.update(id, dto);
    return this.findOne(id);
  }

  async remove(id: number) {
    const item = await this.findOne(id);
    await this.inventoryRepository.remove(item);
    return { message: 'Inventory item deleted' };
  }
}
