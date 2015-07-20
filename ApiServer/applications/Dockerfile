FROM baseimage
# sign name
MAINTAINER zju

#copy the war file into the webapps

RUN mkdir /root/apache-tomcat-7.0.42-v2/webapps/warDir
COPY /newimage  /root/apache-tomcat-7.0.42-v2/webapps/
WORKDIR /root/apache-tomcat-7.0.42-v2/webapps/warDir/

RUN /usr/local/java/jdk1.7.0_67/bin/jar -xvf ../warName

RUN cp ../../lib/apm-trace-0.0.1-SNAPSHOT.jar ./WEB-INF/lib/
RUN cp ../../lib/apm-transaction-0.0.1-SNAPSHOT.jar ./WEB-INF/lib/
# port 22
EXPOSE 22

# Tomcat 8080
EXPOSE 8080

CMD /root/apache-tomcat-7.0.42-v2/bin/start.sh