FROM tomcat:7.0.91-jre8


RUN   apt-get -y install software-properties-common
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys EEA14886
RUN add-apt-repository "deb http://ppa.launchpad.net/webupd8team/java/ubuntu trusty main"
RUN apt-get update && \
    apt-get -y install default-jdk
RUN apt-get install -y curl unzip

WORKDIR  /

RUN mkdir ovitmp && \
    cd ovitmp && \
    curl https://iweb.dl.sourceforge.net/project/dcm4che/Oviyam/2.7.1/Oviyam-2.7.1-bin.zip > oviyam.zip && \
        unzip oviyam.zip
        
RUN rm -R /usr/local/tomcat/webapps/ROOT/
RUN cp /ovitmp/Oviyam-2.7.1-bin/Oviyam-2.7.1-bin/oviyam2.war /usr/local/tomcat/webapps/ROOT.war
RUN cp /ovitmp/Oviyam-2.7.1-bin/tomcat/*.jar  /usr/local/tomcat/lib 
RUN unzip /usr/local/tomcat/webapps/ROOT.war -d  /usr/local/tomcat/webapps/tmp && \
    rm /usr/local/tomcat/webapps/tmp/WEB-INF/web.xml&& \
    rm /usr/local/tomcat/webapps/ROOT.war
COPY web.xml /usr/local/tomcat/webapps/tmp/WEB-INF/web.xml
RUN jar -cvf /usr/local/tomcat/webapps/ROOT.war /usr/local/tomcat/webapps/tmp
COPY tomcat-users.xml /usr/local/tomcat/conf/tomcat-users.xml


