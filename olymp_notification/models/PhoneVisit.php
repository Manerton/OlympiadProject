<?php

namespace app\models;

use Yii;

/**
 * This is the model class for table "phone_visit".
 *
 * @property int $id
 * @property int|null $code
 * @property string $phone_number
 * @property int $type
 * @property string|null $comment
 */
class PhoneVisit extends \yii\db\ActiveRecord
{


    /**
     * {@inheritdoc}
     */
    public static function tableName()
    {
        return 'phone_visit';
    }

    /**
     * {@inheritdoc}
     */
    public function rules()
    {
        return [
            [['code', 'comment'], 'default', 'value' => null],
            [['code', 'type'], 'integer'],
            [['phone_number', 'type'], 'required'],
            [['comment'], 'string'],
            [['phone_number'], 'string', 'max' => 64],
        ];
    }

    /**
     * {@inheritdoc}
     */
    public function attributeLabels()
    {
        return [
            'id' => 'ID',
            'code' => 'Code',
            'phone_number' => 'Phone Number',
            'type' => 'Type',
            'comment' => 'Comment',
        ];
    }

    public static function fill(
        $phoneNumber,
        $type,
        $code = NULL,
        $comment = NULL
    ){
        $entity = new static();
        $entity->phone_number = $phoneNumber;
        $entity->type = $type;
        $entity->code = $code;
        $entity->comment = $comment;
        return $entity;
    }
}