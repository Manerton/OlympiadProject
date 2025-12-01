<?php

namespace app\controllers\api;

use app\components\CheckHelper;
use app\components\CodeHelper;
use app\components\RedisComponent;
use app\services\ApiService;
use app\services\SmsService;
use Yii;
use yii\web\Controller;

class SmsController extends Controller
{
    public SmsService $smsService;
    public function __construct(
        $id,
        $module,
        SmsService $smsService,
        $config = []
    )
    {
        $this->smsService = $smsService;
        parent::__construct($id, $module, $config);
    }

    public function actionSendCode(){
        $data = Yii::$app->request->post();
        $phone = $data['phone'] ?? null;
        $requestToken = $data['requestToken'] ?? null;
        $code = CodeHelper::generateCode();
        if (CheckHelper::checkAccess($requestToken)) {
            $response = $this->smsService->sendSMS($phone, $code);
            RedisComponent::set($phone, $code);
            return Yii::$app->response->data = json_encode([
                'status' => 200,
                'code' => $code,
                'response' => $response
            ]);
        }
        return Yii::$app->response->data = json_encode(['status' => 404]);
    }
}