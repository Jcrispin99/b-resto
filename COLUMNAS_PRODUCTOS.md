# 🍽️ Columnas de Tablas de Productos y Menú

Estructura detallada de las tablas del módulo de productos obtenidas del MCP.

> **⚠️ IMPORTANTE:** Este documento refleja correcciones al MCP original:
> - Se separaron las categorías en `inventory_categories` (inventario) y `product_categories` (menú)
> - Se eliminó `menu_category_id` de `product_template` (redundante)
> - Se identificaron las **tablas base** que deben crearse primero

---

## 🏗️ TABLAS BASE (Crear Primero)

Estas tablas NO tienen dependencias y deben crearse antes que las demás:

### 📦 `inventory_categories` (CUSTOM - Nueva)

**Propósito:** Categorización para el sistema de inventario (materias primas, insumos, etc.).

| Columna | Tipo | Nullable | Default | Descripción |
|---------|------|----------|---------|-------------|
| `id` | bigint | NO | auto | ID único |
| `parent_id` | bigint | YES | null | FK a categoría padre (self-reference) |
| `name` | varchar | NO | - | Nombre de la categoría |
| `full_name` | varchar | YES | null | Nombre completo con jerarquía |
| `is_active` | boolean | NO | true | Estado activo/inactivo |
| `created_at` | timestamp | YES | null | Fecha de creación |
| `updated_at` | timestamp | YES | null | Fecha de actualización |

**Notas:**
- Esta tabla es CUSTOM (no existe en MCP original)
- Se usa para productos almacenables/inventariables
- Ejemplos: "Carnes", "Verduras", "Lácteos", "Bebidas"

---

### 🍕 `product_categories` (Del MCP)

**Propósito:** Categorización para el menú/productos vendibles.

| Columna | Tipo | Nullable | Default | Descripción |
|---------|------|----------|---------|-------------|
| `id` | bigint | NO | auto | ID único |
| `parent_id` | bigint | YES | null | FK a categoría padre (self-reference) |
| `type` | varchar | NO | 'menu' | Tipo de categoría |
| `name` | varchar | NO | - | Nombre de la categoría |
| `full_name` | varchar | YES | null | Nombre completo con jerarquía |
| `is_active` | boolean | NO | true | Estado activo/inactivo |
| `created_at` | timestamp | YES | null | Fecha de creación |
| `updated_at` | timestamp | YES | null | Fecha de actualización |

**Relaciones:**
- Auto-referencia: `parent_id` apunta a otra categoría
- Tiene muchos: `product_template`

**Ejemplos:** "Pizzas", "Hamburguesas", "Bebidas", "Postres"

---

### 🍳 `kitchen_stations` (Del MCP)

**Propósito:** Estaciones de cocina donde se preparan los productos.

| Columna | Tipo | Nullable | Default | Descripción |
|---------|------|----------|---------|-------------|
| `id` | bigint | NO | auto | ID único |
| `branch_id` | bigint | NO | - | FK a sucursal (**NOTA:** crear tabla branches) |
| `name` | varchar | NO | - | Nombre de la estación |
| `description` | varchar | YES | null | Descripción |
| `printer_ip` | varchar | YES | null | IP de impresora de cocina |
| `order` | integer | NO | 0 | Orden de visualización |
| `is_active` | boolean | NO | true | Estado activo/inactivo |
| `created_at` | timestamp | YES | null | Fecha de creación |
| `updated_at` | timestamp | YES | null | Fecha de actualización |

**Relaciones:**
- Pertenece a: `branches` (sucursales)
- Tiene muchos: `product_template`

**Ejemplos:** "Parrilla", "Fríos", "Bebidas", "Hornos"

> ⚠️ **DEPENDENCIA:** Requiere tabla `branches` (sucursales)

---

### 🏪 `warehouses` (Del MCP)

**Propósito:** Almacenes/bodegas donde se guarda inventario.

| Columna | Tipo | Nullable | Default | Descripción |
|---------|------|----------|---------|-------------|
| `id` | bigint | NO | auto | ID único |
| `branch_id` | bigint | NO | - | FK a sucursal |
| `code` | varchar | NO | - | Código del almacén |
| `name` | varchar | NO | - | Nombre del almacén |
| `is_active` | boolean | NO | true | Estado activo/inactivo |
| `created_at` | timestamp | YES | null | Fecha de creación |
| `updated_at` | timestamp | YES | null | Fecha de actualización |

**Relaciones:**
- Pertenece a: `branches`
- Usada por: `inventories`

**Ejemplos:** "Almacén Principal", "Cocina", "Barra"

> ⚠️ **DEPENDENCIA:** Requiere tabla `branches`

---

## 📦 TABLAS DE PRODUCTOS

### 1️⃣ `product_template` (Plantilla de Productos) - **CORREGIDO**

**Propósito:** Maestro de productos, contiene la información común de un producto.

| Columna | Tipo | Nullable | Default | Descripción |
|---------|------|----------|---------|-------------|
| `id` | bigint | NO | auto | ID único |
| `inventory_category_id` | bigint | YES | null | FK a `inventory_categories` (productos almacenables) |
| `category_id` | bigint | NO | - | FK a `product_categories` (categoría de menú) |
| `unit_id` | bigint | NO | - | FK a `units` |
| `name` | varchar | NO | - | Nombre del producto |
| `description` | text | YES | null | Descripción larga |
| `internal_reference` | varchar | YES | null | Código interno |
| `barcode` | varchar | YES | null | Código de barras |
| `product_type` | varchar | NO | 'storable' | Tipo: `storable`, `service`, `consumable` |
| `can_be_sold` | boolean | NO | false | ¿Se puede vender? |
| `can_be_purchased` | boolean | NO | true | ¿Se puede comprar? |
| `can_be_stocked` | boolean | NO | true | ¿Se puede almacenar? |
| `sale_price` | numeric | NO | 0 | Precio de venta |
| `is_active` | boolean | NO | true | Estado activo/inactivo |
| `created_at` | timestamp | YES | null | Fecha de creación |
| `updated_at` | timestamp | YES | null | Fecha de actualización |
| `deleted_at` | timestamp | YES | null | Soft delete |
| `kitchen_station_id` | bigint | YES | null | FK a `kitchen_stations` |

**Relaciones:**
- Pertenece a: `inventory_categories`, `product_categories`, `units`, `kitchen_stations`
- Tiene muchos: `product_product` (variantes)

**Lógica de Categorías:**
- **`inventory_category_id`**: Para productos que afectan inventario (storable, consumable)
  - Ejemplos: Carne de res, Queso, Tomate
- **`category_id`**: Categoría en el menú visible al cliente
  - Ejemplos: Pizzas, Hamburguesas, Bebidas

**Notas:**
- ✅ **CORRECCIÓN:** Se eliminó `menu_category_id` (redundante)
- ✅ **CORRECCIÓN:** Se separó `category_id` en `inventory_category_id` + `category_id`
- Si `product_type = 'storable'` → debe tener `inventory_category_id`

---

## 2️⃣ `product_product` (Variantes de Productos)

**Propósito:** Variantes específicas de un producto template (ej: Pizza Grande, Pizza Mediana).

| Columna | Tipo | Nullable | Default | Descripción |
|---------|------|----------|---------|-------------|
| `id` | bigint | NO | auto | ID único |
| `template_id` | bigint | NO | - | FK a `product_template` |
| `sku` | varchar | NO | - | SKU único de la variante |
| `barcode` | varchar | YES | null | Código de barras específico |
| `sale_price` | numeric | YES | null | Precio de venta (sobrescribe template) |
| `is_active` | boolean | NO | true | Estado activo/inactivo |
| `created_at` | timestamp | YES | null | Fecha de creación |
| `updated_at` | timestamp | YES | null | Fecha de actualización |
| `deleted_at` | timestamp | YES | null | Soft delete |

**Relaciones:**
- Pertenece a: `product_template`
- Muchos a muchos: `product_attribute_values` (via `attribute_value_product`)

---

## 3️⃣ `product_categories` - ✅ Ver sección "Tablas Base"

---

## 4️⃣ `product_attributes` (Atributos de Productos)

**Propósito:** Define tipos de atributos (ej: "Tamaño", "Temperatura", "Extras").

| Columna | Tipo | Nullable | Default | Descripción |
|---------|------|----------|---------|-------------|
| `id` | bigint | NO | auto | ID único |
| `name` | varchar | NO | - | Nombre del atributo (ej: "Tamaño") |
| `created_at` | timestamp | YES | null | Fecha de creación |
| `updated_at` | timestamp | YES | null | Fecha de actualización |

**Relaciones:**
- Tiene muchos: `product_attribute_values`

**Ejemplo:**
```
id=1, name="Tamaño"
id=2, name="Temperatura"
id=3, name="Extras"
```

---

## 5️⃣ `product_attribute_values` (Valores de Atributos)

**Propósito:** Valores específicos de un atributo (ej: "Pequeño", "Mediano", "Grande").

| Columna | Tipo | Nullable | Default | Descripción |
|---------|------|----------|---------|-------------|
| `id` | bigint | NO | auto | ID único |
| `attribute_id` | bigint | NO | - | FK a `product_attributes` |
| `value` | varchar | NO | - | Valor del atributo |
| `created_at` | timestamp | YES | null | Fecha de creación |
| `updated_at` | timestamp | YES | null | Fecha de actualización |

**Relaciones:**
- Pertenece a: `product_attributes`

**Ejemplo:**
```
id=1, attribute_id=1, value="Pequeño"
id=2, attribute_id=1, value="Mediano"
id=3, attribute_id=1, value="Grande"
id=4, attribute_id=2, value="Caliente"
id=5, attribute_id=2, value="Frío"
```

---

## 6️⃣ `product_template_attribute_lines` (Líneas de Atributos por Template)

**Propósito:** Define qué atributos tiene un producto template.

| Columna | Tipo | Nullable | Default | Descripción |
|---------|------|----------|---------|-------------|
| `id` | bigint | NO | auto | ID único |
| `product_template_id` | bigint | NO | - | FK a `product_template` |
| `attribute_id` | bigint | NO | - | FK a `product_attributes` |
| `created_at` | timestamp | YES | null | Fecha de creación |
| `updated_at` | timestamp | YES | null | Fecha de actualización |

**Relaciones:**
- Pertenece a: `product_template`, `product_attributes`

**Ejemplo:**
```
"Pizza" (template_id=1) tiene atributo "Tamaño" (attribute_id=1)
```

---

## 7️⃣ `product_template_attribute_line_values` (Valores por Línea)

**Propósito:** Especifica qué valores de atributo están disponibles para ese template.

| Columna | Tipo | Nullable | Default | Descripción |
|---------|------|----------|---------|-------------|
| `id` | bigint | NO | auto | ID único |
| `product_template_attribute_line_id` | bigint | NO | - | FK a `product_template_attribute_lines` |
| `product_attribute_value_id` | bigint | NO | - | FK a `product_attribute_values` |
| `created_at` | timestamp | YES | null | Fecha de creación |
| `updated_at` | timestamp | YES | null | Fecha de actualización |

**Relaciones:**
- Pertenece a: `product_template_attribute_lines`, `product_attribute_values`

**Ejemplo:**
```
Para "Pizza" con atributo "Tamaño", los valores disponibles son: "Pequeño", "Mediano", "Grande"
```

---

## 8️⃣ `attribute_value_product` (Relación Valores-Productos)

**Propósito:** Tabla pivot que asocia valores de atributos con variantes de productos.

| Columna | Tipo | Nullable | Default | Descripción |
|---------|------|----------|---------|-------------|
| `id` | bigint | NO | auto | ID único |
| `attribute_value_id` | bigint | NO | - | FK a `product_attribute_values` |
| `product_id` | bigint | NO | - | FK a `product_product` |
| `created_at` | timestamp | YES | null | Fecha de creación |
| `updated_at` | timestamp | YES | null | Fecha de actualización |

**Relaciones:**
- Muchos a muchos entre: `product_attribute_values` y `product_product`

**Ejemplo:**
```
product_id=1 (Pizza Grande) tiene attribute_value_id=3 ("Grande")
product_id=2 (Pizza Mediana) tiene attribute_value_id=2 ("Mediano")
```

---

## 9️⃣ `combos` (Combos/Paquetes)

**Propósito:** Paquetes de productos con precio especial.

| Columna | Tipo | Nullable | Default | Descripción |
|---------|------|----------|---------|-------------|
| `id` | bigint | NO | auto | ID único |
| `name` | varchar | NO | - | Nombre del combo |
| `description` | text | YES | null | Descripción del combo |
| `price` | numeric | NO | - | Precio final del combo |
| `regular_price` | numeric | NO | - | Precio regular (suma de productos) |
| `discount_percentage` | numeric | NO | 0 | Porcentaje de descuento |
| `image` | varchar | YES | null | Ruta de la imagen |
| `start_date` | date | NO | - | Fecha de inicio de vigencia |
| `end_date` | date | YES | null | Fecha de fin (null = sin límite) |
| `is_active` | boolean | NO | true | Estado activo/inactivo |
| `created_at` | timestamp | YES | null | Fecha de creación |
| `updated_at` | timestamp | YES | null | Fecha de actualización |
| `deleted_at` | timestamp | YES | null | Soft delete |

**Relaciones:**
- Tiene muchos: `combo_items`

---

## 🔟 `combo_items` (Items de Combos)

**Propósito:** Productos que componen un combo.

| Columna | Tipo | Nullable | Default | Descripción |
|---------|------|----------|---------|-------------|
| `id` | bigint | NO | auto | ID único |
| `combo_id` | bigint | NO | - | FK a `combos` |
| `product_template_id` | bigint | NO | - | FK a `product_template` |
| `quantity` | integer | NO | 1 | Cantidad del producto |
| `allow_substitution` | boolean | NO | false | ¿Permitir sustituciones? |
| `created_at` | timestamp | NO | CURRENT_TIMESTAMP | Fecha de creación |

**Relaciones:**
- Pertenece a: `combos`, `product_template`

---

## 📊 Diagrama de Relaciones

```
product_categories (jerárquica)
    ↓
product_template ─┬─→ product_product (variantes)
    ↓             │       ↓
units             │   attribute_value_product
                  │       ↓
kitchen_stations  │   product_attribute_values
                  │       ↓
                  └──→ product_template_attribute_lines
                          ↓
                      product_template_attribute_line_values
                          
product_attributes
    ↓
product_attribute_values

combos
    ↓
combo_items → product_template
```

---

### ✨ Tablas que Puedes Simplificar

1. **Sistema de Atributos**
   - 5 tablas solo para manejar variantes es complejo
   - Considera usar JSON o solo 3 tablas (ver sugerencias arriba)

2. **Relaciones Polimórficas**
   - `imageables` y `productables` (otras tablas, no en este módulo)
   - En Go, considera usar tablas específicas por tipo o columnas JSON

---

## 📋 Orden de Implementación Actualizado

### Fase 0: Dependencias Previas ⚠️

```
0. branches (sucursales) - CREAR PRIMERO
   └─ Necesaria para: kitchen_stations, warehouses
```

### Fase 1: Tablas Base 🏗️

```
1. inventory_categories (CUSTOM - nueva)
2. product_categories (del MCP)  
3. kitchen_stations (depende de branches)
4. warehouses (depende de branches)
```

### Fase 2: Productos Core 🍽️

```
5. product_template
   └─ Depende de: inventory_categories, product_categories, units, kitchen_stations
   
6. product_product (variantes)
   └─ Depende de: product_template
```

### Fase 3: Atributos (Opcional - Puedes Simplificar) 🔄

```
7. product_attributes
8. product_attribute_values
   └─ Depende de: product_attributes
   
9. product_template_attribute_lines
   └─ Depende de: product_template, product_attributes
   
10. product_template_attribute_line_values
    └─ Depende de: product_template_attribute_lines, product_attribute_values
    
11. attribute_value_product (pivot)
    └─ Depende de: product_attribute_values, product_product
```

### Fase 4: Combos 🎁

```
12. combos
13. combo_items
    └─ Depende de: combos, product_template
```

---

## 📊 Resumen de Correcciones

| Aspecto | MCP Original | Corrección Aplicada |
|---------|-------------|---------------------|
| Categorías | Solo `category_id` | `inventory_category_id` + `category_id` |
| Menu Category | `menu_category_id` | ❌ Eliminado (redundante) |
| Categoría Inventario | No existía | ✅ `inventory_categories` (CUSTOM) |
| Categoría Menú | `product_categories` | ✅ Mantiene mismo nombre |
| Dependencias | No claras | ✅ Identificadas (branches, units) |

---

## ✅ Checklist de Implementación

**Tablas Base:**
- [ ] `branches` (crear primero, no en este doc)
- [ ] `inventory_categories` (CUSTOM)
- [ ] `product_categories`
- [ ] `kitchen_stations`
- [ ] `warehouses`

**Productos:**
- [ ] `product_template`
- [ ] `product_product`

**Atributos (considerar simplificar):**
- [ ] `product_attributes`
- [ ] `product_attribute_values`
- [ ] `product_template_attribute_lines`
- [ ] `product_template_attribute_line_values`
- [ ] `attribute_value_product`

**Combos:**
- [ ] `combos`
- [ ] `combo_items`

---

**Generado:** 2025-12-22  
**Fuente:** MCP Database (PostgreSQL)  
**Correcciones:** Usuario (separación de categorías, eliminación de redundancias)

