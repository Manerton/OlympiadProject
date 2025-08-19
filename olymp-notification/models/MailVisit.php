<?php

namespace app\models;

use Yii;

/**
 * This is the model class for table "mail_visit".
 *
 * @property int $id
 * @property int|null $code
 * @property string $email
 * @property int $type
 * @property string|null $comment
 */
class MailVisit extends \yii\db\ActiveRecord
{


    /**
     * {@inheritdoc}
     */
    public static function tableName()
    {
        return 'mail_visit';
    }

    /**
     * {@inheritdoc}
     */
    public function rules()
    {
        return [
            [['code', 'comment'], 'default', 'value' => null],
            [['code', 'type'], 'integer'],
            [['email', 'type'], 'required'],
            [['comment'], 'string'],
            [['email'], 'string', 'max' => 128],
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
            'email' => 'Email',
            'type' => 'Type',
            'comment' => 'Comment',
        ];
    }

    public static function fill(
        $email,
        $type,
        $code = NULL,
        $comment = NULL
    ){
        $entity = new static();
        $entity->email = $email;
        $entity->type = $type;
        $entity->code = $code;
        $entity->comment = $comment;
        return $entity;
    }
}
