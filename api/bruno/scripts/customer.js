const baseUrl = bru.interpolate("{{baseUrl}}");
const apiVersion = bru.interpolate("{{apiVersion}}");

const createCustomer = async (data) => {
  const cachedCustomer = bru.getVar(data.name)
  if (cachedCustomer) {
    return cachedCustomer
  }

  console.log("📂 Creating customer...")

  try {
    const result = await bru.sendRequest({
      url: `${baseUrl}/api/${apiVersion}/customers`,
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${bru.getVar('jwt_token')}`
      },
      data
    })

    if (!result.data?.success) {
      throw new Error(`Failed to create customer: ${result.data?.error || 'Unknown error'}`)
    }

    bru.setVar(data.name, result.data.data)
    console.log("✅ Customer created successfully")

    return result.data.data
  } catch (error) {
    console.error("❌ Customer creation failed:", error.message)
    throw error
  }
}

const deleteCustomer = async (customerId) => {
  if (!customerId) {
    console.warn("⚠️ No customer ID provided for deletion")
    return
  }

  console.log("🧹 Deleting customer...")

  try {
    await bru.sendRequest({
      url: `${baseUrl}/api/${apiVersion}/customers/${customerId}`,
      method: "DELETE",
      headers: {
        "Authorization": `Bearer ${bru.getVar('jwt_token')}`
      }
    })

    console.log("✅ Customer deleted successfully")
  } catch (error) {
    console.error("❌ Customer deletion failed:", error.message)
    throw error
  }
}

module.exports = {
  createCustomer,
  deleteCustomer
}
