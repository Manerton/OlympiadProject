<?php
use yii\helpers\Html;
use yii\helpers\Url;

$this->title = 'Notification via WebSocket';
?>
<div class="notification">
    <h1><?= Html::encode($this->title) ?></h1>

    <form method = "POST"  id="wsForm" " action="<?= Url::to(['notification/send']) ?>">
        <?= Html::hiddenInput(Yii::$app->request->csrfParam, Yii::$app->request->csrfToken) ?>
        <div class="form-group">
            <input type="text" id="wsMessage" class="form-control" name = "message" placeholder="Введите сообщение" required>
        </div>
        <div class="form-group mt-2">
            <button type="submit" class="btn btn-success btn-lg">Отправить</button>
        </div>
    </form>

    <hr>

    <h3>Лог WebSocket:</h3>
    <div id="wsLog" style="background:#f9f9f9; padding:10px; border:1px solid #ccc; height:200px; overflow:auto;"></div>
</div>

<script>
    /*const socket = new WebSocket("ws://notification.olymp.local/notify");
    const logDiv = document.getElementById("wsLog");
    function log(msg) {
        const p = document.createElement("p");
        p.textContent = msg;
        logDiv.appendChild(p);
        logDiv.scrollTop = logDiv.scrollHeight;
    }
    socket.onopen = () => log("✅ WebSocket соединение установлено");
    socket.onmessage = (event) => log("📩 Пришло сообщение: " + event.data);
    socket.onclose = () => log("❌ WebSocket соединение закрыто");
    socket.onerror = (e) => log("⚠️ Ошибка: " + e.message);
    document.getElementById("wsForm").addEventListener("submit", function (e) {
        e.preventDefault();
        const message = document.getElementById("wsMessage").value;
        if (message && socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify({
                event: "olymp_admin",
                data: message
            }));
            log("➡️ Отправлено: " + message);
            document.getElementById("wsMessage").value = "";
        } else {
            log("⚠️ Нельзя отправить: соединение не открыто");
        }
    });*/
</script>
