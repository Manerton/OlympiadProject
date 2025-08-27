<?php

namespace app\controllers;

use app\components\CheckHelper;
use app\components\CodeHelper;
use app\components\RedisComponent;
use app\models\MailVisit;
use app\repositories\MailVisitRepository;
use app\services\MailService;
use app\services\ApiService;
use components\MessageDictionary;
use Yii;
use yii\web\Controller;

class EmailController extends Controller
{
    private MailService $mailService;
    private MailVisitRepository $mailVisitRepository;
    public function __construct(
        $id,
        $module,
        MailService $mailService,
        MailVisitRepository $mailVisitRepository,
        $config = []
    )
    {
        $this->mailService = $mailService;
        $this->mailVisitRepository = $mailVisitRepository;
        parent::__construct($id, $module, $config);
    }

    public function actionSendCode()
    {
        $data = Yii::$app->request->post(); // весь массив из JSON POST
        $email = $data['email'] ?? null;
        $message = $data['message'] ?? null;
        $requestToken = $data['requestToken'] ?? null;

        $code = CodeHelper::generateCode();
        if (CheckHelper::checkAccess($requestToken)) {
            $this->mailService->send(
                $email,
                'Тема сообщения',
                'code',
                ['code' => $code]
            );
            RedisComponent::set($email, $code);
            $model = MailVisit::fill($email, MessageDictionary::CODE_MESSAGE, $code, 'default code message text');
            $this->mailVisitRepository->save($model);
            return Yii::$app->response->data = [
                'status' => 200,
                'code' => $code
            ];
        }
        return Yii::$app->response->data = ['status' => 404];
    }
    public function beforeAction($action)
    {
        if ($action->id === 'send-code') {
            $this->enableCsrfValidation = false;
        }
        return parent::beforeAction($action);
    }

}