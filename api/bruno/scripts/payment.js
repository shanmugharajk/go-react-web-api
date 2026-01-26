const baseUrl = bru.interpolate("{{baseUrl}}");
const apiVersion = bru.interpolate("{{apiVersion}}");

const createVendorPayment = async (data) => {
  try {
    const result = await bru.sendRequest({
      url: `${baseUrl}/api/${apiVersion}/vendor-payments`,
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${bru.getVar('jwt_token')}`
      },
      data
    })

    if (!result.data?.success) {
      throw new Error(`Failed to create vendor payment: ${result.data?.error || 'Unknown error'}`)
    }

    return result.data.data
  } catch (error) {
    console.error("❌ Vendor payment creation failed:", error.message)
    throw error
  }
}

const getVendorPayment = async (paymentId) => {
  if (!paymentId) {
    console.warn("⚠️ No payment ID provided for lookup")
    return
  }

  try {
    const result = await bru.sendRequest({
      url: `${baseUrl}/api/${apiVersion}/vendor-payments/${paymentId}`,
      method: "GET",
      headers: {
        "Authorization": `Bearer ${bru.getVar('jwt_token')}`
      }
    })

    if (!result.data?.success) {
      throw new Error(`Failed to fetch vendor payment: ${result.data?.error || 'Unknown error'}`)
    }

    return result.data.data
  } catch (error) {
    console.error("❌ Vendor payment fetch failed:", error.message)
    throw error
  }
}

module.exports = {
  createVendorPayment,
  getVendorPayment
}
