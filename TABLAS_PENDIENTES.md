# Análisis de Tablas MCP - Estado de Migración

## 📊 Resumen General

**Total de tablas en MCP:** 55 tablas  
**Tablas implementadas en Go:** 6 tablas  
**Tablas pendientes:** 49 tablas  
**Tablas de Laravel/Sistema:** ~8 tablas (se pueden omitir)

---

## ✅ Tablas Ya Implementadas (6)

| Tabla | Modelo Go | Estado |
|-------|-----------|--------|
| `companies` | `Company` | ✅ Implementado |
| `payment_methods` | `PaymentMethod` | ✅ Implementado |
| `taxes` | `Tax` | ✅ Implementado |
| `units` | `Unit` | ✅ Implementado |
| `users` | `User` | ✅ Implementado |
| - | `Claims` | ✅ Implementado (utilidad JWT) |

---

## 🔴 Tablas Pendientes de Migración (49)

### 🏢 **Configuración y Organización (5 tablas)**

| # | Tabla | Prioridad | Descripción |
|---|-------|-----------|-------------|
| 1 | `partners` | 🔴 Alta | Socios/Proveedores del negocio |
| 2 | `warehouses` | 🔴 Alta | Almacenes/Bodegas |
| 3 | `table_areas` | 🟡 Media | Áreas de mesas (salón, terraza, etc.) |
| 4 | `tables` | 🟡 Media | Mesas del restaurante |
| 5 | `kitchen_stations` | 🟡 Media | Estaciones de cocina (parrilla, fríos, etc.) |

---

### 🍽️ **Productos y Menú (12 tablas)**

| # | Tabla | Prioridad | Descripción |
|---|-------|-----------|-------------|
| 6 | `product_template` | 🔴 Alta | Plantilla de productos (maestro) |
| 7 | `product_product` | 🔴 Alta | Variantes de productos |
| 8 | `product_categories` | 🔴 Alta | Categorías de productos |
| 9 | `product_attributes` | 🟡 Media | Atributos (tamaño, temperatura, etc.) |
| 10 | `product_attribute_values` | 🟡 Media | Valores de atributos (pequeño, grande, etc.) |
| 11 | `product_template_attribute_lines` | 🟡 Media | Líneas de atributos por plantilla |
| 12 | `product_template_attribute_line_values` | 🟡 Media | Valores específicos por línea |
| 13 | `attribute_value_product` | 🟡 Media | Relación muchos-a-muchos atributos-productos |
| 14 | `combos` | 🟢 Baja | Combos/Paquetes |
| 15 | `combo_items` | 🟢 Baja | Items de los combos |
| 16 | `product_menu_settings` | 🟢 Baja | Configuración de menú |
| 17 | `branch_menu_availability` | 🟢 Baja | Disponibilidad por sucursal |

---

### 📦 **Inventario y Stock (4 tablas)**

| # | Tabla | Prioridad | Descripción |
|---|-------|-----------|-------------|
| 18 | `inventories` | 🔴 Alta | Control de inventario |
| 19 | `stock_transfers` | 🟡 Media | Transferencias entre almacenes |
| 20 | `recipes` | 🟡 Media | Recetas (ingredientes de productos) |
| 21 | `purchase_orders` | 🟡 Media | Órdenes de compra |

---

### 🧾 **Órdenes y Ventas (9 tablas)**

| # | Tabla | Prioridad | Descripción |
|---|-------|-----------|-------------|
| 22 | `orders` | 🔴 Alta | Órdenes de venta |
| 23 | `order_items` | 🔴 Alta | Items de las órdenes |
| 24 | `order_payments` | 🔴 Alta | Pagos de las órdenes |
| 25 | `sale_orders` | 🟡 Media | Órdenes de venta adicionales |
| 26 | `reservations` | 🟢 Baja | Reservaciones de mesas |
| 27 | `kitchen_tickets` | 🟡 Media | Tickets de cocina |
| 28 | `kitchen_ticket_items` | 🟡 Media | Items de tickets de cocina |
| 29 | `pos_terminals` | 🟡 Media | Terminales POS |
| 30 | `terminal_journals` | 🟡 Media | Journals por terminal |

---

### 💰 **Caja y Finanzas (3 tablas)**

| # | Tabla | Prioridad | Descripción |
|---|-------|-----------|-------------|
| 31 | `cash_registers` | 🔴 Alta | Cajas registradoras |
| 32 | `cash_movements` | 🔴 Alta | Movimientos de efectivo |
| 33 | `journals` | 🔴 Alta | Diarios contables |
| 34 | `sequences` | 🔴 Alta | Secuencias para numeración |

---

### 🔐 **Permisos y Roles (4 tablas)**

| # | Tabla | Prioridad | Descripción | Notas |
|---|-------|-----------|-------------|-------|
| 35 | `permissions` | 🟡 Media | Permisos del sistema | Considerar simplificar |
| 36 | `roles` | 🟡 Media | Roles de usuario | Simplificar vs Laravel |
| 37 | `role_has_permissions` | 🟡 Media | Relación roles-permisos | Evaluar necesidad |
| 38 | `model_has_permissions` | 🟢 Baja | Permisos por modelo | Probablemente omitir |
| 39 | `model_has_roles` | 🟢 Baja | Roles por modelo | Probablemente omitir |

---

### 🖼️ **Multimedia y Polimórficas (2 tablas)**

| # | Tabla | Prioridad | Descripción | Notas |
|---|-------|-----------|-------------|-------|
| 40 | `imageables` | 🟡 Media | Relación polimórfica de imágenes | Considerar otra estrategia |
| 41 | `productables` | 🟡 Media | Relación polimórfica de productos | Evaluar necesidad |

---

### 🚫 **Tablas de Laravel/Sistema (8 tablas - OMITIR)**

| # | Tabla | Acción | Razón |
|---|-------|--------|-------|
| 42 | `migrations` | ⏭️ Omitir | Sistema de migraciones Laravel |
| 43 | `cache` | ⏭️ Omitir | Cache de Laravel |
| 44 | `cache_locks` | ⏭️ Omitir | Locks de cache Laravel |
| 45 | `sessions` | ⏭️ Omitir | Sesiones Laravel (usar JWT) |
| 46 | `password_reset_tokens` | ⏭️ Omitir | Reset de contraseñas Laravel |
| 47 | `personal_access_tokens` | ⏭️ Omitir | Tokens Laravel (usar JWT) |
| 48 | `failed_jobs` | ⏭️ Omitir | Jobs fallidos Laravel |
| 49 | `jobs` | ⏭️ Omitir | Sistema de jobs Laravel |
| 50 | `job_batches` | ⏭️ Omitir | Batches de jobs Laravel |

---

## 📋 Plan de Implementación Sugerido

### Fase 1: Core del Negocio (Prioridad Alta) 🔴

**Orden sugerido:**

1. **Productos Base**
   - `product_template`
   - `product_product`
   - `product_categories`

2. **Organización**
   - `partners`
   - `warehouses`

3. **Inventario**
   - `inventories`

4. **Órdenes y Ventas**
   - `orders`
   - `order_items`
   - `order_payments`

5. **Caja**
   - `journals`
   - `sequences`
   - `cash_registers`
   - `cash_movements`

### Fase 2: Funcionalidades Intermedias (Prioridad Media) 🟡

6. **Atributos de Productos**
   - `product_attributes`
   - `product_attribute_values`
   - `product_template_attribute_lines`
   - `product_template_attribute_line_values`

7. **Mesas y Cocina**
   - `table_areas`
   - `tables`
   - `kitchen_stations`
   - `kitchen_tickets`
   - `kitchen_ticket_items`

8. **Stock y Compras**
   - `stock_transfers`
   - `recipes`
   - `purchase_orders`

9. **POS**
   - `pos_terminals`
   - `terminal_journals`

10. **Permisos (simplificados)**
    - `permissions`
    - `roles`
    - `role_has_permissions`

### Fase 3: Funcionalidades Adicionales (Prioridad Baja) 🟢

11. **Combos y Menú**
    - `combos`
    - `combo_items`
    - `product_menu_settings`
    - `branch_menu_availability`

12. **Otras**
    - `reservations`
    - `sale_orders`

---

## 💡 Recomendaciones

### ✨ Tablas que Puedes Simplificar

1. **Sistema de Permisos Laravel Spatie**
   - Las tablas `model_has_permissions` y `model_has_roles` son muy específicas de Laravel
   - Considera un sistema más simple de roles en Go

2. **Relaciones Polimórficas**
   - `imageables` y `productables` usan relaciones polimórficas de Laravel
   - En Go, considera usar tablas específicas por tipo o columnas JSON

3. **Tokens y Sesiones**
   - Ya estás usando JWT, no necesitas `sessions` ni `personal_access_tokens`

### 🗑️ Tablas que Definitivamente Puedes Omitir

- Todas las de sistema Laravel (cache, migrations, jobs, etc.)
- `password_reset_tokens` (implementa tu propio sistema)

### 🔄 Tablas que Puedes Rediseñar

1. **`product_template` + `product_product`**
   - Evalúa si realmente necesitas esta separación
   - En muchos casos un solo modelo `Product` con variantes JSON puede ser más simple

2. **Atributos de productos**
   - 5 tablas solo para atributos es complejo
   - Considera usar JSON o un diseño más simple

---

## 📊 Estadísticas Finales

| Categoría | Cantidad |
|-----------|----------|
| **Total tablas MCP** | 55 |
| **Implementadas** | 6 |
| **Omitir (Laravel)** | 8 |
| **Pendientes reales** | 41 |
| **Prioridad Alta** | ~15 |
| **Prioridad Media** | ~18 |
| **Prioridad Baja** | ~8 |

---

## ✅ Próximos Pasos Recomendados

1. **Revisa este documento** y decide qué tablas realmente necesitas
2. **Simplifica** las que puedas (especialmente permisos y polimórficas)
3. **Empieza por Fase 1** (productos, órdenes, caja)
4. **Crea los modelos en Go** siguiendo tu estructura actual
5. **Implementa los controladores** y rutas para cada módulo

---

**Generado el:** 2025-12-22  
**Proyecto:** b-resto (Go Backend)
