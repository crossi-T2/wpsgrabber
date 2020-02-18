%define use_systemd (0%{?fedora} && 0%{?fedora} >= 18) || (0%{?rhel} && 0%{?rhel} >= 7) || (0%{?suse_version} && 0%{?suse_version} >=1210)

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
wpsgrabber is a tool for watching and analyse processing reports from a 52north WPS server (https://52north.org).

%define debug_package %{nil}

%prep

%build


%install
%{__mkdir} -p %{buildroot}%{_sysconfdir}/%{name}
%{__install} -m644 %{_sourcedir}/configs/config.yaml %{buildroot}%{_sysconfdir}/%{name}
%{__mkdir} -p %{buildroot}/usr/local/bin/
cp %{_sourcedir}/%{name} %{buildroot}/usr/local/bin/
%if %{use_systemd}
# systemd-specific files
%{__mkdir} -p %{buildroot}%{_unitdir}
%{__install} -m644 %{_sourcedir}/init/%{name}.service \
    %{buildroot}%{_unitdir}/%{name}.service
%else
# SYSV init files
%{__install} -D -m 755 %{_sourcedir}/init/%{name} %{buildroot}%{_initrddir}/%{name}
%endif

%post

%postun

%clean
rm -rf %{buildroot}

%files
%config(noreplace) %{_sysconfdir}/%{name}
%config(noreplace) %{_sysconfdir}/%{name}/*
/usr/local/bin/%{name}
%if %{use_systemd}
%{_unitdir}/%{name}.service
%else
%{_initrddir}/%{name}
%endif

%changelog
