<?php

namespace app\services;

use Symfony\Component\Mailer\Transport;
use Symfony\Component\Mime\Email;
use yii\base\Component;
use Symfony\Component\Mailer\Mailer as SymfonyMailerCore;
class MailSymfonyService extends Component
{
    /** @var SymfonyMailerCore */
    private $mailer;

    /** @var string */
    public $dsn;

    public function init()
    {
        parent::init();
        $transport = Transport::fromDsn($this->dsn);
        $this->mailer = new SymfonyMailerCore($transport);
    }

    public function compose($view = null, array $params = [])
    {
        return new Email();
    }

    public function send(Email $message): bool
    {
        $this->mailer->send($message);
        return true;
    }
}