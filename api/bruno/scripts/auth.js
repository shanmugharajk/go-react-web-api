const baseUrl = bru.interpolate("{{baseUrl}}");
const apiVersion = bru.interpolate("{{apiVersion}}");

const login = async () => {
  if (bru.getVar('jwt_token')) {
    return
  }

  console.log("🔐 Logging in test user...")

  const testUserEmail = bru.interpolate("{{testUserEmail}}");
  const testUserPassword = bru.interpolate("{{testUserPassword}}");

  const response = await bru.sendRequest({
    url: `${baseUrl}/api/${apiVersion}/auth/token/login`,
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    data: {
      email: testUserEmail,
      password: testUserPassword,
    }
  });
  
  bru.setVar('jwt_token', response.data.data.accessToken)
  console.log("✅ User logged in successfully")
}

const register = async () => {
  if (bru.getVar('jwt_token')) {
    return
  }

  console.log("📝 Registering test user...")

  const testUserEmail = bru.interpolate("{{testUserEmail}}");
  const testUserPassword = bru.interpolate("{{testUserPassword}}");
  const testUserName = bru.interpolate("{{testUserName}}");

  const response = await bru.sendRequest({
    url: `${baseUrl}/api/${apiVersion}/auth/token/register`,
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    data: {
      email: testUserEmail,
      password: testUserPassword,
      name: testUserName
    }
  })

  bru.setVar('jwt_token', response.data.data.accessToken)
  console.log("✅ User registered successfully")
}

module.exports = {
  login,
  register
}