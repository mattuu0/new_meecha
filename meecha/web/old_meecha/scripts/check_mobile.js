var ua = navigator.userAgent;
if (!(ua.indexOf('iPhone') > 0 || ua.indexOf('iPad') > 0 || ua.indexOf('Android') > 0 || ua.indexOf('Mobile') > 0 )) {
    window.location.href = window.location.origin + '/meecha/desktop_view';
}

