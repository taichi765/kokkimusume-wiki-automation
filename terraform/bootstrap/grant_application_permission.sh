APPLICATION_READWRITE_OWNEDBY_ID=18a4783c-866b-4cc7-a460-3d5e5662c884
APP_ID=f09650f5-9059-46d8-910f-619a2738737a
PRINCIPAL_ID=831f36dd-9254-4fe3-a2cf-752edee59d94

az rest --method GET --url "https://graph.microsoft.com/v1.0/servicePrincipals?\
$filter=appId eq $APP_ID&\
$select=id,displayName,appId,appRoles"

MSGRAPH_RESOURCE_ID=96f140a1-f138-40e7-a79b-a32767835731

BODY="{\"resourceId\":\"${MSGRAPH_RESOURCE_ID}\",\"principalId\":\"${PRINCIPAL_ID}\",\"appRoleId\":\"${APPLICATION_READWRITE_OWNEDBY_ID}\"}"

az rest -m POST -u "https://graph.microsoft.com/v1.0/servicePrincipals/${PRINCIPAL_ID}/appRoleAssignedTo" -b $BODY