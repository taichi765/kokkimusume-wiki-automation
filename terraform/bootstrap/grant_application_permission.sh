APPLICATION_READWRITE_OWNEDBY_ID=18a4783c-866b-4cc7-a460-3d5e5662c884
MSGRAPH_APP_ID=00000003-0000-0000-c000-000000000000
PRINCIPAL_ID=$(terraform output -raw service_principal_id)
echo "PRINCIPAL_ID: ${PRINCIPAL_ID}"

echo "getting MSGRAPH_RESOURCE_ID..."
MSGRAPH_RESOURCE_ID=$(az rest --method GET --url "https://graph.microsoft.com/v1.0/servicePrincipals?\
\$filter=appId eq '${MSGRAPH_APP_ID}'&\
\$select=id,displayName,appId,appRoles"\
|jq -r ".value\
|map(select(.displayName==\"Microsoft Graph\"))\
|if length == 1 then .[0] \
elif length == 0 then error (\"can't find Microsoft Graph service principal\") \
else error(\"too many service principal with name Microsoft Graph\")
end\
|.id")
echo "MSGRAPH_RESOURCE_ID: ${MSGRAPH_RESOURCE_ID}"

BODY="{\"resourceId\":\"${MSGRAPH_RESOURCE_ID}\",\"principalId\":\"${PRINCIPAL_ID}\",\"appRoleId\":\"${APPLICATION_READWRITE_OWNEDBY_ID}\"}"
echo "request body: ${BODY}"

az rest -m POST -u "https://graph.microsoft.com/v1.0/servicePrincipals/${PRINCIPAL_ID}/appRoleAssignedTo" -b $BODY
