<?php

use yii\db\Migration;

class m250606_185010_init_tables extends Migration
{
    /**
     * {@inheritdoc}
     */
    public function safeUp()
    {
        // PhoneVisit table
        $this->createTable('phone_visit', [
            'id' => $this->primaryKey(),
            'code' => $this->integer()->null(),
            'phone_number' => $this->string(64)->notNull(),
            'type' => $this->integer()->notNull(),
            'comment' => $this->text()->null(),
        ]);

        // MailVisit table
        $this->createTable('mail_visit', [
            'id' => $this->primaryKey(),
            'code' => $this->integer()->null(),
            'email' => $this->string(128)->notNull(),
            'type' => $this->integer()->notNull(),
            'comment' => $this->text()->null(),
        ]);
    }

    /**
     * {@inheritdoc}
     */
    public function safeDown()
    {
        $this->dropTable('phone_visit');
        $this->dropTable('mail_visit');
    }

    /*
    // Use up()/down() to run migration code without a transaction.
    public function up()
    {

    }

    public function down()
    {
        echo "m250606_185010_init_tables cannot be reverted.\n";

        return false;
    }
    */
}
