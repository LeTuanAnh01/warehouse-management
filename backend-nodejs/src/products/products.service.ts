import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Product } from './product.entity';

@Injectable()
export class ProductsService {
  constructor(
    @InjectRepository(Product)
    private productsRepository: Repository<Product>,
  ) {}

  findAll() {
    return this.productsRepository.find();
  }

  async findOne(id: number) {
    const product = await this.productsRepository.findOne({ where: { id } });
    if (!product) throw new NotFoundException('Product not found');
    return product;
  }

  create(dto: Partial<Product>) {
    const product = this.productsRepository.create(dto);
    return this.productsRepository.save(product);
  }

  async update(id: number, dto: Partial<Product>) {
    await this.productsRepository.update(id, dto);
    return this.findOne(id);
  }

  async remove(id: number) {
    const product = await this.findOne(id);
    await this.productsRepository.remove(product);
    return { message: 'Product deleted' };
  }
}
