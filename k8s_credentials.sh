kubectl create secret docker-registry regcred --docker-server=https://index.docker.io/v2/ --docker-username=<username>--docker-password=XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX --docker-email=<emailid>

kubectl create secret generic hanacloud --from-literal=driverName='hdb' --from-literal=hdbDsn='xyz'

kubectl create secret generic cpi --from-literal=cpi_client_id='<client id>' --from-literal=cpi_client_secret='<client secret>' --from-literal=cpi_token_endpoint='https://681769b2trial.authentication.us10.hana.ondemand.com/oauth/token' --from-literal=cpi_api_endpoint='https://681769b2trial.it-cpitrial05.cfapps.us10-001.hana.ondemand.com'



