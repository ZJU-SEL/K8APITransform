# VERSION 0.0.1
# using a tomcat7 image
FROM jre7-tomcat7
# sign name
MAINTAINER zju

#copy the war file into the webapps
ADD /applications/10.10.103.250/logtest-0.0.1.war_deploy/logtest.war  /usr/local/tomcat/webapps/


# port 22
EXPOSE 22

# Tomcat 8080
EXPOSE 8080
