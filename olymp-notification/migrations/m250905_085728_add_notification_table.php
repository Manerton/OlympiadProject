<?php

use yii\db\Migration;

class m250905_085728_add_notification_table extends Migration
{
    /**
     * {@inheritdoc}
     */
    public function safeUp()
    {
        $this->createTable('notification', [
            'id' => $this->primaryKey(),
            'user_id' => $this->string(64)->notNull(),
            'message' => $this->string(1000)->notNull(),
            'status' => $this->integer(),
        ]);
    }

    /**
     * {@inheritdoc}
     */
    public function safeDown()
    {
        $this->dropTable('notification');
        return true;
    }

    /*
    // Use up()/down() to run migration code without a transaction.
    public function up()
    {

    }

    public function down()
    {
        echo "m250905_085728_add_notification_table cannot be reverted.\n";

        return false;
    }
    */
}
