# Bruno API Tests

This directory contains API tests for the POS system using [Bruno](https://www.usebruno.com/).

## Directory Structure

```
bruno/
├── auth/                    # Authentication endpoints
├── system/                  # System health checks
├── entities/                # Core entity CRUD tests
│   ├── customers/          # Customer entity tests
│   ├── products/           # Product entity tests
│   ├── product-categories/ # Product category tests
│   └── folder.bru          # Shared authentication setup
├── scripts/                 # Shared helper functions
│   ├── auth.js             # Authentication helpers
│   ├── customer.js         # Customer test helpers
│   ├── product.js          # Product test helpers
│   └── utils.js            # Shared utilities (UUID validation, etc.)
└── environments/           # Environment configurations
    └── local.bru           # Local development environment
```

## Test Organization

### Entities vs Journeys

- **`entities/`** - Tests for individual CRUD operations on core entities (customers, products, etc.)
- **`journeys/`** (planned) - End-to-end workflow tests (sales transactions, stock management, etc.)

### Folder-level Fixtures

Each test folder may contain a `folder.bru` file that runs before every test in that folder:

- **`entities/folder.bru`** - Registers test user and obtains JWT token
- **`entities/products/folder.bru`** - Creates shared product category and product (cached across tests)

These fixtures use caching to avoid redundant API calls and persist for the entire test run.

## Naming Conventions

### Test Data Naming

Test data follows a consistent naming pattern to avoid collisions:

```javascript
// Pattern: scope.entity.operation[.detail]
"entities.customer.create"
"entities.product.folder.productCategory"
"entities.product.delete"
```

### Variable Naming

Bruno variables use the same namespaced pattern:

```javascript
// Good
bru.setVar('entities.customer.create.id', customerId)
bru.setVar('entities.product.folder.productCategoryId', categoryId)

// Avoid generic names that could collide
// Bad: bru.setVar('testId', id)
```

## Running Tests

### Prerequisites

1. Start the API server (see root Makefile)
2. Ensure database is in a clean state

### Commands

```bash
# From project root
make test              # Run all entity tests
make test-watch        # Run tests in watch mode

# Using Bruno CLI directly
cd api/bruno
bru run --env local --tags entities              # All entity tests
bru run --env local --tags customers             # Customer tests only
bru run --env local --tags products              # Product tests only
```

### Test Workflow

1. **Terminal 1:** `make run` - Starts server with fresh database
2. **Terminal 2:** `make test` - Runs Bruno tests

## Writing New Tests

### Individual Entity Tests

For CRUD operations on a single entity:

```javascript
// 1. Create test data in pre-request
script:pre-request {
  const helper = require('./scripts/helper.js')

  const data = {
    name: "entities.entity.operation",
    // ... other fields
  }

  const result = await helper.createEntity(data)
  bru.setVar('entities.entity.operation.id', result.id.toString())
}

// 2. Clean up in post-response
script:post-response {
  const helper = require('./scripts/helper.js')
  const id = bru.getVar('entities.entity.operation.id')

  if (id) {
    await helper.deleteEntity(id)
  }
}
```

### Using Folder Fixtures

Tests can reference shared folder fixtures:

```javascript
// Use the folder-level product category
body:json {
  {
    "name": "Test Product",
    "categoryId": "{{entities.product.folder.productCategoryId}}"
  }
}
```

### Test Assertions

Use the shared UUID validator:

```javascript
test("should return valid UUID", function() {
  const { isValidUUID } = require('./scripts/utils')
  const body = res.getBody()
  expect(isValidUUID(body.data.id)).to.be.true
})
```

## Helper Functions

### Available Helpers

All helper functions include:
- ✅ Error handling with descriptive messages
- ✅ Null/undefined guards
- ✅ Logging with emoji indicators
- ✅ Caching support (where applicable)

**`scripts/auth.js`**
- `register()` - Register test user and get JWT token
- `login()` - Login existing user

**`scripts/customer.js`**
- `createCustomer(data)` - Create customer (cached by name)
- `deleteCustomer(id)` - Delete customer

**`scripts/product.js`**
- `createProductCategory(data)` - Create category (cached by name)
- `createProduct(data)` - Create product (cached by name)
- `deleteProduct(id)` - Delete product
- `deleteProductCategory(id)` - Delete category

**`scripts/utils.js`**
- `isValidUUID(str)` - Validate UUID format
- `uuidRegex` - UUID regex pattern (prefer `isValidUUID()`)

### Caching Mechanism

Helper functions cache by the `name` field to avoid duplicate creations:

```javascript
// First call: Makes API request
const category1 = await createProductCategory({
  name: "entities.product.folder.productCategory"
})

// Second call with same name: Returns cached result (no API call)
const category2 = await createProductCategory({
  name: "entities.product.folder.productCategory"
})

// category1 === category2 (same object reference)
```

## Best Practices

### Test Isolation

- ✅ Each test should clean up its own data
- ✅ Use namespaced variable names to avoid collisions
- ✅ Don't rely on execution order between tests
- ✅ Folder fixtures are acceptable for shared read-only data

### Test Data

- ✅ Use descriptive, namespaced names
- ✅ Keep test data minimal but realistic
- ✅ Clean up test data in `post-response` hooks
- ⚠️ Folder fixtures persist for the entire test run (cleaned on next run)

### Performance

- ✅ Use folder fixtures for shared data (products, categories)
- ✅ Leverage caching in helper functions
- ✅ Run tests in parallel where possible (use tags)

### Maintainability

- ✅ Use helper functions instead of inline API calls
- ✅ Follow consistent naming conventions
- ✅ Add comments for complex test scenarios
- ✅ Keep assertions focused and clear

## Database Management

Tests use a SQLite database that should be reset before each test run:

```bash
# Manual reset
rm -rf api/data/pos.db

# Automated (via Makefile)
make run  # Automatically resets database
```

In CI/CD, create a fresh database for each pipeline run.

## Future Enhancements

### Journey Tests (Planned)

For complex workflows spanning multiple entities:

```
journeys/
├── sales/
│   ├── complete-sale.bru       # End-to-end sale workflow
│   └── sale-with-discount.bru  # Sale with promotions
├── inventory/
│   └── stock-adjustment.bru    # Stock management workflow
└── folder.bru                   # Journey-level fixtures
```

### Test Data Builders

**⚠️ IMPORTANT: Builders are for JOURNEY tests only, NOT entity tests!**

Test data builders (`scripts/builders.js`) provide a fluent interface for creating complex test scenarios. Use them ONLY for journey tests where you need many variations of test data.

**When to use builders:**
- ✅ Journey tests with complex workflows (sales, inventory, multi-step processes)
- ✅ Tests requiring many objects with different configurations
- ✅ Scenarios with multiple variations (VIP customer vs. regular customer)

**When NOT to use builders:**
- ❌ Entity CRUD tests (customers, products, categories)
- ❌ Simple single-operation tests
- ❌ Tests with straightforward data requirements

**Entity tests (current) - Keep it simple:**
```javascript
// ✅ GOOD - Simple inline data for entity tests
body:json {
  {
    "name": "entities.customer.create",
    "email": "customer@example.com",
    "balance": 100.50
  }
}
```

**Journey tests (future) - Use builders:**
```javascript
// ✅ GOOD - Builders for complex journey tests
script:pre-request {
  const { CustomerBuilder, ProductBuilder, Presets } = require('./scripts/builders.js')

  // Complex scenario: VIP customer purchasing multiple products
  const vipCustomer = new CustomerBuilder()
    .withName("VIP Customer")
    .withBalance(10000)
    .build()

  const expensiveProduct = Presets.product.expensive(categoryId)
  const cheapProduct = Presets.product.cheap(categoryId)

  // Create sale with multiple items...
}
```

**Available builders:**
- `CustomerBuilder` - For customer test data
- `ProductBuilder` - For product test data
- `ProductCategoryBuilder` - For category test data
- `Presets` - Common scenarios (premium customer, expensive product, etc.)

See `scripts/builders.js` for full documentation and examples.

## Troubleshooting

### Tests Failing with 404 Errors

- Ensure the server is running (`make run`)
- Check that folder fixtures completed successfully
- Verify database has been reset

### Tests Failing with 401 Unauthorized

- Check that `entities/folder.bru` ran successfully
- Verify JWT token is being set in folder pre-request
- Ensure auth endpoints are working

### Caching Issues

- Caching is based on the `name` field in test data
- Clear cache by restarting the test run
- Each test run starts with a fresh cache

### Cleanup Failures

- Check server logs for FK constraint violations
- Ensure deletion order respects foreign keys (products before categories)
- Verify IDs are being stored correctly in variables

## Resources

- [Bruno Documentation](https://docs.usebruno.com/)
- [Bruno CLI Reference](https://docs.usebruno.com/bru-cli/overview)
- [Bruno Scripting Guide](https://docs.usebruno.com/testing/script/javascript-reference)
