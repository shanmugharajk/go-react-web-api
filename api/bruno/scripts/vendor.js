const baseUrl = bru.interpolate("{{baseUrl}}");
const apiVersion = bru.interpolate("{{apiVersion}}");

const createVendor = async (data) => {
  const cachedVendor = bru.getVar(data.name)
  if (cachedVendor) {
    return cachedVendor
  }

  try {
    const result = await bru.sendRequest({
      url: `${baseUrl}/api/${apiVersion}/vendors`,
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${bru.getVar('jwt_token')}`
      },
      data
    })

    if (!result.data?.success) {
      throw new Error(`Failed to create vendor: ${result.data?.error || 'Unknown error'}`)
    }

    bru.setVar(data.name, result.data.data)

    return result.data.data
  } catch (error) {
    console.error("❌ Vendor creation failed:", error.message)
    throw error
  }
}

const deleteVendor = async (vendorId) => {
  if (!vendorId) {
    console.warn("⚠️ No vendor ID provided for deletion")
    return
  }

  try {
    await bru.sendRequest({
      url: `${baseUrl}/api/${apiVersion}/vendors/${vendorId}`,
      method: "DELETE",
      headers: {
        "Authorization": `Bearer ${bru.getVar('jwt_token')}`
      }
    })

    bru.deleteVar(vendorId)
  } catch (error) {
    console.error("❌ Vendor deletion failed:", error.message)
    throw error
  }
}

const getVendor = async (vendorId) => {
  if (!vendorId) {
    console.warn("⚠️ No vendor ID provided for lookup")
    return
  }

  try {
    const result = await bru.sendRequest({
      url: `${baseUrl}/api/${apiVersion}/vendors/${vendorId}`,
      method: "GET",
      headers: {
        "Authorization": `Bearer ${bru.getVar('jwt_token')}`
      }
    })

    if (!result.data?.success) {
      throw new Error(`Failed to fetch vendor: ${result.data?.error || 'Unknown error'}`)
    }

    return result.data.data
  } catch (error) {
    console.error("❌ Vendor fetch failed:", error.message)
    throw error
  }
}

module.exports = {
  createVendor,
  deleteVendor,
  getVendor
}
