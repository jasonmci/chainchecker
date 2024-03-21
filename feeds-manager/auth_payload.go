package main

// GeneratePayload generates the GraphQL payload with dynamic username and password
func GenerateAuthPayload(username, password string) map[string]interface{} {
    return map[string]interface{}{
        "operationName": "Login",
        "variables": map[string]interface{}{
            "input": map[string]interface{}{
                "email":    username,
                "password": password,
            },
        },
        "query": `mutation Login($input: LoginInput!) {
            login(input: $input) {
                session {
                    id
                    token
                    expiresAt
                    revoked
                    permits
                    __typename
                }
                errors {
                    ...MutationErrorFields
                    __typename
                }
                __typename
            }
        }
        
        fragment MutationErrorFields on MutationError {
            ...GenericErrorFields
            ...InvalidInputErrorFields
            ...UniqueViolationErrorFields
            __typename
        }
        
        fragment GenericErrorFields on GenericError {
            message
            __typename
        }
        
        fragment InvalidInputErrorFields on InvalidInputError {
            inputError: message
            path
            __typename
        }
        
        fragment UniqueViolationErrorFields on UniqueViolationError {
            message
            __typename
        }`,
    }
}
