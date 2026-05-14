import { Controller, Get, Post, Patch, Delete, Param, Body, UseGuards } from '@nestjs/common';
import { WarehousesService } from './warehouses.service';
import { JwtAuthGuard } from '../common/guards/jwt-auth.guard';

@Controller('api/warehouses')
@UseGuards(JwtAuthGuard)
export class WarehousesController {
  constructor(private warehousesService: WarehousesService) {}

  @Get()
  findAll() {
    return this.warehousesService.findAll();
  }

  @Get(':id')
  findOne(@Param('id') id: string) {
    return this.warehousesService.findOne(+id);
  }

  @Post()
  create(@Body() dto: any) {
    return this.warehousesService.create(dto);
  }

  @Patch(':id')
  update(@Param('id') id: string, @Body() dto: any) {
    return this.warehousesService.update(+id, dto);
  }

  @Delete(':id')
  remove(@Param('id') id: string) {
    return this.warehousesService.remove(+id);
  }
}
