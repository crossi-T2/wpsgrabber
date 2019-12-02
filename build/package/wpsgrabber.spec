Name:           wpsgrabber
Url:            https://github.com/crossi-T2/wpsgrabber
License:        AGPLv3
Version:        0.1
Release:        %{_release}
Summary:        wpsgrabber
BuildArch:      noarch
AutoReqProv:    no
BuildRequires:  libtool


%description
wpsgrabber is a tool for watching processing reports from a 52north WPS server (https://52north.org) and producing simple CSV files for further analysis.

%define debug_package %{nil}

%prep

%build


%install
mkdir -p %{buildroot}%{_sysconfdir}/wpsgrabber/init
cp %{_sourcedir}/init/wpsgrabber.service %{buildroot}%{_sysconfdir}/wpsgrabber/init/
cp %{_sourcedir}/init/wpsgrabber %{buildroot}%{_sysconfdir}/wpsgrabber/init/
cp -r %{_sourcedir}/configs/config.yaml %{buildroot}%{_sysconfdir}/wpsgrabber
mkdir -p %{buildroot}/usr/local/bin/
cp %{_sourcedir}/wpsgrabber %{buildroot}/usr/local/bin/

%post
#!/bin/bash
mkdir -p /var/log/wpsgrabber/

centos=$(rpm -E %{rhel})

if [ "${centos}" == "6" ]; then
    cp /etc/wpsgrabber/init/wpsgrabber /etc/init.d/
    chmod +x /etc/init.d/wpsgrabber
else
    if [ "${centos}" == "7" ]; then
        /etc/wpsgrabber/init/wpsgrabber.service /usr/lib/systemd/system/
    fi
fi

%postun
#!/bin/bash
rm -rf /var/log/wpsgrabber/

centos=$(rpm -E %{rhel})

if [ "${centos}" == "6" ]; then
    rm -f /etc/init.d/wpsgrabber
else
    if [ "${centos}" == "7" ]; then
        rm -f /usr/lib/systemd/system/wpsgrabber.service
    fi
fi

%clean
rm -rf %{buildroot}

%files
%config(noreplace) %{_sysconfdir}/*
/usr/local/bin/wpsgrabber

%changelog
