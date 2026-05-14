import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { WarehousesService } from './warehouses.service';
import { WarehousesController } from './warehouses.controller';
import { Warehouse } from './warehouse.entity';

@Module({
  imports: [TypeOrmModule.forFeature([Warehouse])],
  controllers: [WarehousesController],
  providers: [WarehousesService],
})
export class WarehousesModule {}
