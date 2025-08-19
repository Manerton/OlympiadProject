<?php

namespace app\services;

use Yii;
use yii\base\Component;

class MailService extends Component
{
    public $fromEmail;

    public function init()
    {
        parent::init();
        if ($this->fromEmail === null) {
            $this->fromEmail = Yii::$app->params['adminEmail'];
        }
    }

    /**
     * Отправка письма
     *
     * @param string|array $to
     * @param string $subject
     * @param string $view // путь к представлению
     * @param array $params // параметры для view
     * @return bool
     */
    public function send($to, $subject, $view, $params = [])
    {
        return Yii::$app->mailer->compose($view, $params)
            ->setFrom($this->fromEmail)
            ->setTo($to)
            ->setSubject($subject)
            ->send();
    }
}