var modalForm, modalAlert;
var currentIdx, currentDevice;

function showFormNew() {
    const element = document.getElementById("modalTitle");
    element.innerText = "Add new wake-up"
    fillForm("", "", "", "insert");
    modalForm = new bootstrap.Modal(document.getElementById('wolForm'), {backdrop: 'static'});
    modalForm.show()
}

function showFormEdit(idx) {
    const data = document.getElementById("idx-" + idx)
    const element = document.getElementById("modalTitle");
    element.innerText = "Edit wake-up " + data.dataset.device
    fillForm(data.dataset.device, data.dataset.mac, data.dataset.ip, "update")
    modalForm = new bootstrap.Modal(document.getElementById('wolForm'), {backdrop: 'static'});
    modalForm.show()
}

function fillForm(dev, m, i, s) {
    const device = document.getElementById("device-input");
    device.value = dev
    const odevice = document.getElementById("odevice-input");
    odevice.value = dev
    const mac = document.getElementById("mac-input");
    mac.value = m
    const ip = document.getElementById("ip-input");
    ip.value = i
    const scope = document.getElementById("scope-input");
    scope.value = s

    device.focus();
}

function submitForm() {
    var formElement = document.getElementById("form-input");
    formElement.submit();
    modalForm.hide();
}

function deleteAlert(idx) {
    currentIdx = idx
    const data = document.getElementById("idx-" + idx)
    const title = document.getElementById("alertTitle");
    title.innerText = "Delete Wake-Up"
    const message = document.getElementById("alertMessage");
    message.innerText = "Do you want to delete " + data.dataset.device + "?";
    currentDevice = data.dataset.device;
    modalAlert = new bootstrap.Modal(document.getElementById('alertForm'), {backdrop: 'static'});
    modalAlert.show()
}

function deleteItem() {
    window.location.href = "/" + escape(currentDevice) +  "/delete"
    modalAlert.hide()
    console.log("Delete item " + currentIdx)
}

function cloneItem(idx) {
    const data = document.getElementById("idx-" + idx)
    window.location.href = "/" + escape(data.dataset.device) +  "/clone"
}

function qrCode(idx) {
    const data = document.getElementById("idx-" + idx)
    const title = document.getElementById("qrcodeTitle");
    title.innerText = "QR-Code for " + data.dataset.device
    const image = document.getElementById("qrcodeImage")
    image.src = "/" + escape(data.dataset.device) +  "/qrcode"
    modalAlert = new bootstrap.Modal(document.getElementById('qrcodeForm'), {backdrop: 'static'});
    modalAlert.show()
}
