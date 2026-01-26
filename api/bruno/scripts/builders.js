/**
 * Test Data Builders
 *
 * These builders provide a fluent interface for creating test data.
 * Use them in journey tests where you need complex, varied test scenarios.
 *
 * Example:
 *   const customer = new CustomerBuilder()
 *     .withName("John Doe")
 *     .withBalance(1000)
 *     .build()
 */

class CustomerBuilder {
  constructor() {
    this.data = {
      name: "Test Customer",
      email: "test@example.com",
      mobile: "0000000000",
      balance: 0,
      active: true
    }
  }

  withName(name) {
    this.data.name = name
    return this
  }

  withEmail(email) {
    this.data.email = email
    return this
  }

  withMobile(mobile) {
    this.data.mobile = mobile
    return this
  }

  withBalance(balance) {
    this.data.balance = balance
    return this
  }

  inactive() {
    this.data.active = false
    return this
  }

  build() {
    return { ...this.data }
  }
}

class ProductCategoryBuilder {
  constructor() {
    this.data = {
      name: "Test Category",
      description: "A test category"
    }
  }

  withName(name) {
    this.data.name = name
    return this
  }

  withDescription(description) {
    this.data.description = description
    return this
  }

  build() {
    return { ...this.data }
  }
}

class ProductBuilder {
  constructor() {
    this.data = {
      name: "Test Product",
      description: "A test product",
      price: 0,
      stock: 0,
      isActive: true
    }
  }

  withName(name) {
    this.data.name = name
    return this
  }

  withDescription(description) {
    this.data.description = description
    return this
  }

  withPrice(price) {
    this.data.price = price
    return this
  }

  withStock(stock) {
    this.data.stock = stock
    return this
  }

  inCategory(categoryId) {
    this.data.categoryId = categoryId
    return this
  }

  inactive() {
    this.data.isActive = false
    return this
  }

  build() {
    return { ...this.data }
  }
}

/**
 * Preset builders for common test scenarios
 */
const Presets = {
  customer: {
    withBalance: (amount) => new CustomerBuilder()
      .withName("Customer with Balance")
      .withEmail("balance@example.com")
      .withBalance(amount)
      .build(),

    premium: () => new CustomerBuilder()
      .withName("Premium Customer")
      .withEmail("premium@example.com")
      .withBalance(10000)
      .build(),

    new: () => new CustomerBuilder()
      .withName("New Customer")
      .withEmail("new@example.com")
      .withBalance(0)
      .build()
  },

  product: {
    expensive: (categoryId) => new ProductBuilder()
      .withName("Premium Product")
      .withPrice(999.99)
      .withStock(10)
      .inCategory(categoryId)
      .build(),

    cheap: (categoryId) => new ProductBuilder()
      .withName("Budget Product")
      .withPrice(9.99)
      .withStock(100)
      .inCategory(categoryId)
      .build(),

    outOfStock: (categoryId) => new ProductBuilder()
      .withName("Out of Stock Product")
      .withPrice(49.99)
      .withStock(0)
      .inCategory(categoryId)
      .build()
  }
}

module.exports = {
  CustomerBuilder,
  ProductCategoryBuilder,
  ProductBuilder,
  Presets
}
