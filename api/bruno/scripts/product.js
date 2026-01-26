const baseUrl = bru.interpolate("{{baseUrl}}");
const apiVersion = bru.interpolate("{{apiVersion}}");

const createProductCategory = async (data) => {
  const cachedCategory = bru.getVar(data.name)
  if (cachedCategory) {
    return cachedCategory
  }

  try {
    const result = await bru.sendRequest({
      url: `${baseUrl}/api/${apiVersion}/products/categories`,
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${bru.getVar('jwt_token')}`
      },
      data
    })

    if (!result.data?.success) {
      throw new Error(`Failed to create category: ${result.data?.error || 'Unknown error'}`)
    }

    bru.setVar(data.name, result.data.data)

    return result.data.data
  } catch (error) {
    console.error("❌ Product category creation failed:", error.message)
    throw error
  }
}

const createProduct = async (data) => {
  const cachedProduct = bru.getVar(data.name)
  if (cachedProduct) {
    return cachedProduct
  }

  try {
    const result = await bru.sendRequest({
      url: `${baseUrl}/api/${apiVersion}/products`,
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${bru.getVar('jwt_token')}`
      },
      data
    })

    if (!result.data?.success) {
      throw new Error(`Failed to create product: ${result.data?.error || 'Unknown error'}`)
    }

    bru.setVar(data.name, result.data.data)

    return result.data.data
  } catch (error) {
    console.error("❌ Product creation failed:", error.message)
    throw error
  }
}

const deleteProduct = async (productId) => {
  if (!productId) {
    console.warn("⚠️ No product ID provided for deletion")
    return
  }

  try {
    await bru.sendRequest({
      url: `${baseUrl}/api/${apiVersion}/products/${productId}`,
      method: "DELETE",
      headers: {
        "Authorization": `Bearer ${bru.getVar('jwt_token')}`
      }
    })

    bru.deleteVar(productId)
  } catch (error) {
    console.error("❌ Product deletion failed:", error.message)
    throw error
  }
}

const deleteProductCategory = async (categoryId) => {
  if (!categoryId) {
    console.warn("⚠️ No category ID provided for deletion")
    return
  }

  try {
    await bru.sendRequest({
      url: `${baseUrl}/api/${apiVersion}/products/categories/${categoryId}`,
      method: "DELETE",
      headers: {
        "Authorization": `Bearer ${bru.getVar('jwt_token')}`
      }
    })

    bru.deleteVar(categoryId)
  } catch (error) {
    console.error("❌ Product category deletion failed:", error.message)
    throw error
  }
}

module.exports = {
  createProductCategory,
  createProduct,
  deleteProduct,
  deleteProductCategory
}