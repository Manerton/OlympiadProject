<?php

namespace app\controllers;

use app\components\CheckHelper;
use app\components\CodeHelper;
use app\models\PhoneVisit;
use app\services\PhoneNumberService;
use components\MessageDictionary;
use repositories\PhoneVisitRepository;
use Yii;
use yii\web\Controller;

class PhoneNumberController extends Controller
{
    private PhoneVisitRepository $phoneVisitRepository;
    public function __construct(
        $id,
        $module,
        PhoneVisitRepository $phoneVisitRepository,
        $config = []
    )
    {
        $this->phoneVisitRepository = $phoneVisitRepository;
        parent::__construct($id, $module, $config);
    }

    public function actionSendTelegramCode()
    {
        $code = CodeHelper::generateCode();
        $requestToken = Yii::$app->request->headers->get('requestToken');
        $phoneNumber =  Yii::$app->request->headers->get('phone_number');
        if (CheckHelper::checkAccess($requestToken)) {
            $model = PhoneVisit::fill($phoneNumber, MessageDictionary::CODE_MESSAGE, $code, 'default code message text');
            $this->phoneVisitRepository->save($model);
            //место для вызова PhoneNumberService
            return \Yii::$app->response->data = json_encode(['status' => 200, 'code' => $code]);
        }
        return \Yii::$app->response->data = json_encode(['status' => 404]);
    }
    public function actionSendWhatsappCode()
    {
        $code = CodeHelper::generateCode();
        $requestToken = Yii::$app->request->headers->get('requestToken');
        $phoneNumber =  Yii::$app->request->headers->get('phone_number');
        if (CheckHelper::checkAccess($requestToken)) {
            $model = PhoneVisit::fill($phoneNumber, MessageDictionary::CODE_MESSAGE, $code, 'default code message text');
            $this->phoneVisitRepository->save($model);
            //место для вызова PhoneNumberService
            return \Yii::$app->response->data = json_encode(['status' => 200, 'code' => $code]);
        }
        return \Yii::$app->response->data = json_encode(['status' => 404]);
    }
}