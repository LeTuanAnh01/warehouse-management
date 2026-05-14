import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Warehouse } from './warehouse.entity';

@Injectable()
export class WarehousesService {
  constructor(
    @InjectRepository(Warehouse)
    private warehousesRepository: Repository<Warehouse>,
  ) {}

  findAll() {
    return this.warehousesRepository.find();
  }

  async findOne(id: number) {
    const warehouse = await this.warehousesRepository.findOne({ where: { id } });
    if (!warehouse) throw new NotFoundException('Warehouse not found');
    return warehouse;
  }

  create(dto: Partial<Warehouse>) {
    const warehouse = this.warehousesRepository.create(dto);
    return this.warehousesRepository.save(warehouse);
  }

  async update(id: number, dto: Partial<Warehouse>) {
    await this.warehousesRepository.update(id, dto);
    return this.findOne(id);
  }

  async remove(id: number) {
    const warehouse = await this.findOne(id);
    await this.warehousesRepository.remove(warehouse);
    return { message: 'Warehouse deleted' };
  }
}
