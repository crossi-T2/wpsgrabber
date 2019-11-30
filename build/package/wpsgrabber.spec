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
mkdir -p %{buildroot}%{_unitdir}
cp %{_sourcedir}/init/wpsgrabber.service %{buildroot}%{_unitdir}
mkdir -p %{buildroot}%{_sysconfdir}/wpsgrabber
cp -r %{_sourcedir}/configs/config.yaml %{buildroot}%{_sysconfdir}/wpsgrabber
mkdir -p %{buildroot}/usr/local/bin/
cp %{_sourcedir}/wpsgrabber %{buildroot}/usr/local/bin/

%post
mkdir -p /var/log/wpsgrabber/

%postun

%clean
rm -rf %{buildroot}

%files
%{_unitdir}/*
%config(noreplace) %{_sysconfdir}/*
/usr/local/bin/wpsgrabber

%changelog
