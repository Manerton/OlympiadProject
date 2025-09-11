<?php

namespace app\controllers;

use app\components\CheckHelper;
use app\components\CodeHelper;
use app\components\RedisComponent;
use app\models\MailVisit;
use app\repositories\MailVisitRepository;
use app\components\MessageDictionary;
use Symfony\Component\Mime\Email;
use Yii;
use yii\web\Controller;

class EmailController extends Controller
{
    private MailVisitRepository $mailVisitRepository;
    public function __construct(
        $id,
        $module,
        MailVisitRepository $mailVisitRepository,
        $config = []
    )
    {
        $this->mailVisitRepository = $mailVisitRepository;
        parent::__construct($id, $module, $config);
    }
    public function actionSendCode()
    {
        $data = Yii::$app->request->post();
        $email = $data['email'] ?? null;
        $requestToken = $data['requestToken'] ?? null;
        $code = CodeHelper::generateCode();
        if (CheckHelper::checkAccess($requestToken)) {
            $mail = (new Email())
                ->from(Yii::$app->params['adminEmail'])
                ->to($email)
                ->subject('Письмо с кодом доступа')
                ->text('Код: ')
                ->html($code);
            Yii::$app->mailer->send($mail);
            RedisComponent::set($email, $code);
            $model = MailVisit::fill($email, MessageDictionary::CODE_MESSAGE, $code, 'default code message text');
            $this->mailVisitRepository->save($model);
            return Yii::$app->response->data = json_encode([
                'status' => 200,
                'code' => $code
            ]);
        }
        return Yii::$app->response->data = json_encode(['status' => 404]);
    }
    public function actionSendMessage()
    {
        $data = Yii::$app->request->post();
        $email = $data['email'] ?? null;
        $message = $data['message'] ?? null;
        $requestToken = $data['requestToken'] ?? null;
        if (CheckHelper::checkAccess($requestToken)) {
            $mail = (new Email())
                ->from(Yii::$app->params['adminEmail'])
                ->to($email)
                ->subject('Письмо. ВСоШ')
                ->text('Код: ')
                ->html($message);
            Yii::$app->mailer->send($mail);
            $model = MailVisit::fill($email, MessageDictionary::TEXT_MESSAGE, NULL, $message);
            $this->mailVisitRepository->save($model);
            return Yii::$app->response->data = json_encode([
                'status' => 200,
            ]);
        }
        return Yii::$app->response->data = json_encode(['status' => 404]);
    }
    public function beforeAction($action)
    {
        if ($action->id === 'send-code' || $action->id === 'send-message') {
            $this->enableCsrfValidation = false;
        }
        return parent::beforeAction($action);
    }

}