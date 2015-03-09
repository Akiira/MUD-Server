-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

-- -----------------------------------------------------
-- Schema MUD-Database
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema MUD-Database
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `MUD-Database` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci ;
USE `MUD-Database` ;

-- -----------------------------------------------------
-- Table `MUD-Database`.`Login`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `MUD-Database`.`Login` ;

CREATE TABLE IF NOT EXISTS `MUD-Database`.`Login` (
  `CharacterName` VARCHAR(16) NOT NULL,
  `Password` VARCHAR(16) NOT NULL,
  `WorldLocation` VARCHAR(6) NOT NULL,
  PRIMARY KEY (`CharacterName`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `MUD-Database`.`Character`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `MUD-Database`.`Character` ;

CREATE TABLE IF NOT EXISTS `MUD-Database`.`Character` (
  `CharacterName` VARCHAR(16) NOT NULL,
  `Race` VARCHAR(45) NOT NULL,
  `Strength` INT NOT NULL,
  `Intelligence` INT NOT NULL,
  `Constitution` INT NOT NULL,
  `Dexterity` INT NOT NULL,
  `Charisma` INT NOT NULL,
  `Wisdom` INT NOT NULL,
  PRIMARY KEY (`CharacterName`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `MUD-Database`.`Inventory`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `MUD-Database`.`Inventory` ;

CREATE TABLE IF NOT EXISTS `MUD-Database`.`Inventory` (
  `CharacterName` VARCHAR(16) NOT NULL,
  `ItemID` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`CharacterName`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `MUD-Database`.`Weapons`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `MUD-Database`.`Weapons` ;

CREATE TABLE IF NOT EXISTS `MUD-Database`.`Weapons` (
  `idWeapons` INT NOT NULL,
  `Name` VARCHAR(45) NULL,
  `weaponType` INT NULL,
  PRIMARY KEY (`idWeapons`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `MUD-Database`.`Armour`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `MUD-Database`.`Armour` ;

CREATE TABLE IF NOT EXISTS `MUD-Database`.`Armour` (
  `idArmour` INT NOT NULL,
  `Name` VARCHAR(45) NULL,
  `armourClass` INT NULL,
  `wearLocation` INT NULL,
  PRIMARY KEY (`idArmour`))
ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
